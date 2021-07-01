// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "github.com/Urethramancer/signor/opt"

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

	return nil
}
