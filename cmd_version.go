// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"github.com/grimdork/climate/arg"
)

func cmdVersion(opts *arg.Options) error {
	opt := arg.New("awsec version", "Show version information.")
	opt.SetDefaultHelp(true)
	// No positional args; just parse any remaining args for help support.
	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	pr("awsec version %s (%s)", version, date)
	return nil
}
