// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/opt"
)

// RenameCmd options.
type RenameCmd struct {
	opt.DefaultHelp
	Key string `placeholder:"KEY" help:"Old key name."`
	New string `placeholder:"NAME" help:"New name for the key. Contents will be unchanged."`
}

// Run rename.
func (cmd *RenameCmd) Run(in []string) error {
	if cmd.Help || cmd.Key == "" {
		return opt.ErrUsage
	}

	client, err := getClient()
	if err != nil {
		return nil
	}

	param, err := client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(validKey(cmd.Key)),
		WithDecryption: true,
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
		Values: []string{validKey(cmd.Key)},
	}

	dpi := ssm.DescribeParametersInput{
		MaxResults:       25,
		ParameterFilters: []types.ParameterStringFilter{filter},
	}

	desc, err := client.DescribeParameters(context.Background(), &dpi)
	if err != nil {
		return err
	}

	ppi := &ssm.PutParameterInput{
		Name:        aws.String(validKey(cmd.New)),
		Value:       param.Parameter.Value,
		Type:        param.Parameter.Type,
		Description: desc.Parameters[0].Description,
		Overwrite:   true,
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
		ResourceId:   aws.String(validKey(cmd.New)),
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
