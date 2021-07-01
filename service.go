// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func getClient() (*ssm.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return ssm.NewFromConfig(cfg), nil
}
