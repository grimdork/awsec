package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/climate/arg"
)

func cmdRename(opts *arg.Options) error {
	opt := arg.New("awsec rename", "Rename key.")
	opt.SetDefaultHelp(true)
	opt.SetPositional("KEY", "Old key name.", "", true, arg.VarString)
	opt.SetPositional("NAME", "New name for the key. Contents will be unchanged.", "", true, arg.VarString)

	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	param, err := client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(validKey(opt.GetPosString("KEY"))),
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		return err
	}

	if !askString("Are you sure you want to rename this parameter?") {
		return nil
	}

	filter := types.ParameterStringFilter{
		Key:    aws.String("Name"),
		Option: aws.String("Equals"),
		Values: []string{validKey(opt.GetPosString("KEY"))},
	}

	dpi := ssm.DescribeParametersInput{
		MaxResults:       aws.Int32(25),
		ParameterFilters: []types.ParameterStringFilter{filter},
	}

	desc, err := client.DescribeParameters(context.Background(), &dpi)
	if err != nil {
		return err
	}

	ppi := &ssm.PutParameterInput{
		Name:        aws.String(validKey(opt.GetPosString("NAME"))),
		Value:       param.Parameter.Value,
		Type:        param.Parameter.Type,
		Description: desc.Parameters[0].Description,
		Overwrite:   aws.Bool(true),
	}
	_, err = client.PutParameter(context.Background(), ppi)
	if err != nil {
		return err
	}

	tags, err := client.ListTagsForResource(context.Background(), &ssm.ListTagsForResourceInput{
		ResourceId:   param.Parameter.Name,
		ResourceType: types.ResourceTypeForTaggingParameter,
	})
	if err != nil {
		return err
	}

	_, err = client.AddTagsToResource(context.Background(), &ssm.AddTagsToResourceInput{
		ResourceId:   aws.String(validKey(opt.GetPosString("NAME"))),
		ResourceType: types.ResourceTypeForTaggingParameter,
		Tags:         tags.TagList,
	})
	if err != nil {
		return err
	}

	_, err = client.DeleteParameter(context.Background(), &ssm.DeleteParameterInput{
		Name: param.Parameter.Name,
	})
	return err
}
