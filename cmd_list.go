package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/grimdork/climate/arg"
)

func cmdList(opts *arg.Options) error {
	opt := arg.New("awsec list", "List keys.")
	opt.SetDefaultHelp(true)
	opt.SetOption(arg.GroupDefault, "f", "filter", "Lists all keys starting with this.", "/", false, arg.VarString, nil)
	opt.SetOption(arg.GroupDefault, "o", "output", "Output format.", "table", false, arg.VarString, []any{"table", "compact", "json"})
	err := opt.Parse(opts.Args)
	if err != nil {
		return err
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	filter := types.ParameterStringFilter{
		Key:    aws.String("Name"),
		Option: aws.String("BeginsWith"),
		Values: []string{opt.GetString("filter")},
	}
	input := &ssm.DescribeParametersInput{
		MaxResults:       aws.Int32(50),
		ParameterFilters: []types.ParameterStringFilter{filter},
	}

	w := &tabwriter.Writer{}
	jsonout := jsonSecrets{}
	output := opt.GetString("output")
	switch output {
	case "json":
	case "compact":
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "Secret,Last modified")
	default:
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)
		fmt.Fprintln(w, "Secret\tLast modified\tDescription")
	}

	loop := true
	for loop {
		po, err := client.DescribeParameters(context.Background(), input)
		if err != nil {
			return err
		}

		for _, p := range po.Parameters {
			if p.Description == nil {
				p.Description = aws.String("")
			}

			s := fmt.Sprintf("%s\t%s\t%s\n", *p.Name, p.LastModifiedDate.String(), *p.Description)
			switch output {
			case "json":
				e := jsonSecret{
					Name:        *p.Name,
					LastMod:     p.LastModifiedDate.String(),
					Description: *p.Description,
				}

				jsonout.Secrets = append(jsonout.Secrets, e)
			case "compact":
				t := p.LastModifiedDate
				date := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
				fmt.Fprintf(w, "%s,%s\n", *p.Name, date)
			default:
				fmt.Fprint(w, s)
			}
		}

		if po.NextToken == nil {
			loop = false
		} else {
			input.NextToken = po.NextToken
		}
	}

	switch output {
	case "json":
		data, err := json.MarshalIndent(jsonout, "", "\t")
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", data)
	default:
		w.Flush()
	}

	return nil
}

type jsonSecrets struct {
	Secrets []jsonSecret `json:"secrets"`
}

type jsonSecret struct {
	Name        string `json:"secret"`
	LastMod     string `json:"last_modified"`
	Description string `json:"description,omitempty"`
	Contents    string `json:"contents,omitempty"`
}
