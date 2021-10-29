// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

// VersionCmd options.
type VersionCmd struct{}

// Run tag.
func (cmd *VersionCmd) Run(in []string) error {
	pr("awsec version %s (%s)", version, date)
	return nil
}
