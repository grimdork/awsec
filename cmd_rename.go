// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "github.com/Urethramancer/signor/opt"

// RenameCmd options.
type RenameCmd struct {
	opt.DefaultHelp
	Key string `placeholder:"KEY" help:"Old key name."`
	New string `placeholder:"NAME" help:"New name for the key. Contents will be unchanged."`
}

// Run rename.
func (cmd *RenameCmd) Run(in []string) error {
	if cmd.Help || cmd.Key == "" {
		return opt.ErrUsage
	}

	return nil
}
