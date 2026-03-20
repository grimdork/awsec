package main

import (
	"github.com/grimdork/climate/arg"
)

func cmdVersion(opts *arg.Options) error {
	pr("awsec version %s (%s)", version, date)
	return nil
}
