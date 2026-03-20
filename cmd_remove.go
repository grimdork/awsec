package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/grimdork/climate/arg"
)

func cmdRemove(opts *arg.Options) error {
	opt := arg.New("awsec remove", "Remove key.")
	opt.SetDefaultHelp(true)
	opt.SetPositional("KEY", "Key to remove.", "", true, arg.VarString)
	opt.SetFlag(arg.GroupDefault, "f", "force", "Don't ask for verification.")

	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	if !opt.GetBool("force") {
		if !askString("Are you sure you want to delete this key?") {
			return nil
		}
	}

	_, err = client.DeleteParameter(context.Background(), &ssm.DeleteParameterInput{
		Name: aws.String(opt.GetPosString("KEY")),
	})
	return err
}
