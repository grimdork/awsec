// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Urethramancer/signor/opt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// ListCmd options.
type ListCmd struct {
	opt.DefaultHelp
	Filter string `placeholder:"FILTER" help:"Lists all keys starting with this." default:"/"`
	Desc   bool   `short:"d" help:"Include description."`
}

func (cmd *ListCmd) Run(in []string) error {
	if cmd.Help {
		return opt.ErrUsage
	}

	client, err := getClient()
	if err != nil {
		return nil
	}

	filter := types.ParameterStringFilter{
		Key:    aws.String("Name"),
		Option: aws.String("BeginsWith"),
		Values: []string{validKey(cmd.Filter)},
	}
	input := &ssm.DescribeParametersInput{
		MaxResults:       25,
		ParameterFilters: []types.ParameterStringFilter{filter},
	}

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 0, 8, 1, '\t', 0)
	if cmd.Desc {
		fmt.Fprintln(w, "Secret\tLast modified\tDescription")
	} else {
		fmt.Fprintln(w, "Secret\tLast modified")
	}

	loop := true
	for loop {
		po, err := client.DescribeParameters(context.Background(), input)
		if err != nil {
			return err
		}

		for _, p := range po.Parameters {
			if cmd.Desc {
				if p.Description == nil {
					p.Description = aws.String("")
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", *p.Name, p.LastModifiedDate.String(), *p.Description)
			} else {
				fmt.Fprintf(w, "%s\t%s\n", *p.Name, p.LastModifiedDate.String())
			}
		}

		if po.NextToken == nil {
			loop = false
		} else {
			input.NextToken = po.NextToken
		}
	}
	w.Flush()

	return nil
}
