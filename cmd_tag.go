// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import "github.com/Urethramancer/signor/opt"

// TagCmd options.
type TagCmd struct {
	opt.DefaultHelp
	Key    string `placeholder:"KEY" help:"Key to tag."`
	Remove bool   `short:"r" help:"Remove tags from the key."`
}

// Run tag.
func (cmd *TagCmd) Run(in []string) error {
	if cmd.Help || cmd.Key == "" {
		return opt.ErrUsage
	}

	return nil
}
