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
	Key    string `placeholder:"KEY" help:"Key to set. Any existing key with that name will be replaced."`
	Value  string `placeholder:"VALUE" help:"Value for the key."`
	Desc   string `short:"d" placeholder:"DESCRIPTION" help:"Also set a description for the key."`
	List   bool   `short:"l" help:"The value is a comma-separated StringList."`
	Secure bool   `short:"s" help:"Enable encryption for this key. It will be stored as a regular String."`
}

func (cmd *SetCmd) Run(in []string) error {
	if cmd.Help || cmd.Value == "" {
		return opt.ErrUsage
	}

	svc, err := getClient()
	if err != nil {
		return nil
	}

	ppi := &ssm.PutParameterInput{
		Name:  aws.String(validKey(cmd.Key)),
		Value: aws.String(cmd.Value),
		// Once the type is set, this is what a key will always be. It must be removed to change.
		Type: types.ParameterTypeString,
		// We always want to overwrite, for convenience. Better lock down those permissions.
		Overwrite: true,
	}
	if cmd.Desc != "" {
		ppi.Description = aws.String(cmd.Desc)
	}

	if cmd.List {
		ppi.Type = types.ParameterTypeStringList
	}

	if cmd.Secure {
		ppi.Type = types.ParameterTypeSecureString
	}

	if cmd.List || cmd.Secure {
		// If we change the type to a list or more secure storage, we must delete the old key.
		_, _ = svc.DeleteParameter(context.Background(), &ssm.DeleteParameterInput{Name: aws.String(validKey(cmd.Key))})
		// We don't care if it fails - usually means it didn't exist, and any major failures are caught in the next call.
	}

	_, err = svc.PutParameter(context.Background(), ppi)
	return err
}
