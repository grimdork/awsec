// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"errors"
	"strings"

	"github.com/Urethramancer/signor/opt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// TagCmd options.
type TagCmd struct {
	opt.DefaultHelp
	Key    string   `placeholder:"KEY" help:"Key to tag."`
	Tags   []string `placeholder:"TAG=VALUE" help:"Tag/value pair."`
	Remove bool     `short:"r" help:"Remove the named tags from the key (the values are just keys in this form)."`
}

// Run tag.
func (cmd *TagCmd) Run(in []string) error {
	if cmd.Help || cmd.Key == "" {
		return opt.ErrUsage
	}

	client, err := getClient()
	if err != nil {
		return nil
	}

	if cmd.Remove {
		_, err = client.RemoveTagsFromResource(context.Background(), &ssm.RemoveTagsFromResourceInput{
			ResourceId:   aws.String(cmd.Key),
			ResourceType: types.ResourceTypeForTaggingParameter,
			TagKeys:      cmd.Tags,
		})
		return err
	}

	list := []types.Tag{}
	for _, x := range cmd.Tags {
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
		ResourceId:   aws.String(validKey(cmd.Key)),
		ResourceType: types.ResourceTypeForTaggingParameter,
		Tags:         list,
	})
	return err
}
