package main

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/climate/arg"
)

func cmdTag(opts *arg.Options) error {
	opt := arg.New("awsec tag", "Add tags to a key.")
	opt.SetDefaultHelp(true)
	opt.SetPositional("KEY", "Key to tag.", "", true, arg.VarString)
	opt.SetPositional("TAG=VALUE", "Tag/value pair.", nil, false, arg.VarStringSlice)
	opt.SetFlag(arg.GroupDefault, "r", "remove", "Remove the named tags from the key (the values are just keys in this form).")

	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if opt.GetBool("remove") {
		_, err = client.RemoveTagsFromResource(context.Background(), &ssm.RemoveTagsFromResourceInput{
			ResourceId:   aws.String(opt.GetPosString("KEY")),
			ResourceType: types.ResourceTypeForTaggingParameter,
			TagKeys:      opt.GetPosStringSlice("TAG=VALUE"),
		})
		return err
	}

	list := []types.Tag{}
	for _, x := range opt.GetPosStringSlice("TAG=VALUE") {
		a := strings.Split(x, "=")
		if len(a) != 2 {
			continue
		}

		t := types.Tag{
			Key:   aws.String(a[0]),
			Value: aws.String(a[1]),
		}
		list = append(list, t)
	}

	if len(list) == 0 {
		return errors.New("no valid tags - aborting")
	}

	_, err = client.AddTagsToResource(context.Background(), &ssm.AddTagsToResourceInput{
		ResourceId:   aws.String(validKey(opt.GetPosString("KEY"))),
		ResourceType: types.ResourceTypeForTaggingParameter,
		Tags:         list,
	})
	return err
}
