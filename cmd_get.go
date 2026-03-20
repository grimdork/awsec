package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/climate/arg"
)

func cmdGet(opts *arg.Options) error {
	opt := arg.New("awsec get", "Get value(s) for a key.")
	opt.SetDefaultHelp(true)
	opt.SetPositional("KEY", "Key to fetch.", "", true, arg.VarString)
	opt.SetFlag(arg.GroupDefault, "t", "tags", "Get tags too.")

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

	if opt.GetBool("t") {
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
