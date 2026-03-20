package main

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/climate/arg"
)

func cmdSet(opts *arg.Options) error {
	opt := arg.New("awsec set", "Set the value for a key.")
	opt.SetDefaultHelp(true)
	opt.SetPositional("KEY", "Key to set. Any existing key with that name will be replaced.", "", true, arg.VarString)
	opt.SetPositional("VALUE", "Value for the key.", "", true, arg.VarString)
	opt.SetOption(arg.GroupDefault, "d", "desc", "Also set a description for the key.", "", false, arg.VarString, nil)
	opt.SetFlag(arg.GroupDefault, "l", "list", "The value is a comma-separated StringList.")
	opt.SetFlag(arg.GroupDefault, "s", "secure", "Enable encryption (SecureString) for this key.")
	opt.SetFlag(arg.GroupDefault, "f", "file", "The value is a filename to import as the contents for the key value. Max 4096 bytes.")

	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	key := validKey(opt.GetPosString("KEY"))
	value := opt.GetPosString("VALUE")

	ppi := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(value),
		Type:      types.ParameterTypeString,
		Overwrite: aws.Bool(true),
	}

	desc := opt.GetString("desc")
	if desc != "" {
		ppi.Description = aws.String(desc)
	}

	if opt.GetBool("list") {
		ppi.Type = types.ParameterTypeStringList
	}

	if opt.GetBool("secure") {
		ppi.Type = types.ParameterTypeSecureString
	}

	if opt.GetBool("file") {
		data, err := os.ReadFile(value)
		if err != nil {
			return err
		}

		ppi.Value = aws.String(string(data))

		fi, err := os.Stat(value)
		if err != nil {
			return err
		}

		if fi.Size() > 8192 {
			return errors.New("file size is too large (>8192 bytes)")
		}

		if fi.Size() > 4096 {
			ppi.Tier = types.ParameterTierAdvanced
		}
	}

	if opt.GetBool("list") || opt.GetBool("secure") {
		// If we change the type to a list or more secure storage, we must delete the old key.
		_, _ = client.DeleteParameter(context.Background(), &ssm.DeleteParameterInput{Name: aws.String(key)})
		// We don't care if it fails - usually means it didn't exist, and any major failures are caught in the next call.
	}

	_, err = client.PutParameter(context.Background(), ppi)
	return err
}
