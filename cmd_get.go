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
)

// GetCmd options.
type GetCmd struct {
	opt.DefaultHelp
	Key string `placeholder:"KEY" help:"Key to fetch."`
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
		Name:           aws.String(cmd.Key),
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

	return nil
}
