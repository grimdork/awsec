// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/climate/arg"
)

func cmdBackup(opts *arg.Options) error {
	opt := arg.New("awsec backup", "Back up all keys to S3.")
	opt.SetDefaultHelp(true)
	opt.SetPositional("BUCKET", "S3 bucket to back up everything to.", "", true, arg.VarString)

	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	filter := types.ParameterStringFilter{
		Key:    aws.String("Name"),
		Option: aws.String("BeginsWith"),
		Values: []string{validKey("/")},
	}
	input := &ssm.DescribeParametersInput{
		MaxResults:       aws.Int32(25),
		ParameterFilters: []types.ParameterStringFilter{filter},
	}

	loop := true
	for loop {
		po, err := client.DescribeParameters(context.Background(), input)
		if err != nil {
			return err
		}

		for _, p := range po.Parameters {
			if p.Description == nil {
				p.Description = aws.String("")
			}
			fmt.Printf("%s\t%s\t%s\n", *p.Name, p.LastModifiedDate.String(), *p.Description)
		}

		if po.NextToken == nil {
			loop = false
		} else {
			input.NextToken = po.NextToken
		}
	}

	return nil
}
