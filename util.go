// Copyright (c) 2021 Ronny Bangsund
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"context"
	"os"
	"strings"

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

func validKey(k string) string {
	if k[0] != '/' {
		return "/" + k
	}

	return k
}

func askString(q string) bool {
	r := bufio.NewReader(os.Stdin)
	print(q + " [y/N] ")
	res, _ := r.ReadString('\n')
	res = strings.ReplaceAll(res, "\n", "")
	if res == "y" || res == "Y" {
		return true
	}

	return false
}
