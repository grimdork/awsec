// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/Urethramancer/signor/opt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// GetCmd options.
type GetCmd struct {
	opt.DefaultHelp
	Key  string `placeholder:"KEY" help:"Key to fetch."`
	Tags bool   `short:"t" help:"Get tags too."`
}

func (cmd *GetCmd) Run(in []string) error {
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

	switch param.Parameter.Type {
	case "StringList":
		a := strings.Split(*param.Parameter.Value, ",")
		if len(a) < 2 {
			pr("%s", *param.Parameter.Value)
			return nil
		}
		nl := false
		for _, x := range a {
			fmt.Printf("%s", x)
			if nl {
				pr("")
			} else {
				fmt.Print(" = ")
			}
			nl = !nl
		}

	default:
		pr("%s", *param.Parameter.Value)
	}

	if cmd.Tags {
		tags, err := client.ListTagsForResource(context.Background(), &ssm.ListTagsForResourceInput{
			ResourceId:   param.Parameter.Name,
			ResourceType: types.ResourceTypeForTaggingParameter,
		})
		if err != nil {
			return err
		}

		if len(tags.TagList) == 0 {
			return nil
		}

		pr("Tags:")
		for _, t := range tags.TagList {
			pr("\t%s = %s", *t.Key, *t.Value)
		}
	}
	return nil
}
