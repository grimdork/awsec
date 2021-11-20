// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/grimdork/opt"
)

// RemoveCmd options.
type RemoveCmd struct {
	opt.DefaultHelp
	Key   string `placeholder:"KEY" help:"Key to remove."`
	Force bool   `short:"f" help:"Don't ask for verification."`
}

// Run remove.
func (cmd *RemoveCmd) Run(in []string) error {
	if cmd.Help || cmd.Key == "" {
		return opt.ErrUsage
	}

	client, err := getClient()
	if err != nil {
		return nil
	}

	if !askString("Are you sure you want to delete this key?") {
		return nil
	}

	_, err = client.DeleteParameter(context.Background(), &ssm.DeleteParameterInput{
		Name: aws.String(cmd.Key),
	})
	return err

}
