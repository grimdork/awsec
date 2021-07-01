// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"github.com/Urethramancer/signor/opt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// SetCmd options.
type SetCmd struct {
	opt.DefaultHelp
	Key   string `placeholder:"KEY" help:"Key to set."`
	Value string `placeholder:"VALUE" help:"Value for the key."`
	Desc  string `short:"d" placeholder:"DESCRIPTION" help:"Also set a description for the key."`
	List  bool   `short:"l" help:"The value is a comma-separated StringList."`
}

func (cmd *SetCmd) Run(in []string) error {
	if cmd.Help || cmd.Value == "" {
		return opt.ErrUsage
	}

	svc, err := getClient()
	if err != nil {
		return nil
	}

	out, err := svc.PutParameter(context.Background(), &ssm.PutParameterInput{
		Name:        aws.String(cmd.Key),
		Value:       aws.String(cmd.Value),
		Description: aws.String(""),
		Type:        types.ParameterTypeString,
	})
	if err != nil {
		return err
	}

	pr("%+v", out)
	return nil
}
