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
	"github.com/grimdork/opt"
)

// BackupCmd options.
type BackupCmd struct {
	opt.DefaultHelp
	Bucket string `placeholder:"BUCKET" help:"S3 bucket to back up everything to."`
}

func (cmd *BackupCmd) Run(in []string) error {
	if cmd.Help || cmd.Bucket == "" {
		return opt.ErrUsage
	}

	client, err := getClient()
	if err != nil {
		return nil
	}

	filter := types.ParameterStringFilter{
		Key:    aws.String("Name"),
		Option: aws.String("BeginsWith"),
		Values: []string{validKey("/")},
	}
	input := &ssm.DescribeParametersInput{
		MaxResults:       25,
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
