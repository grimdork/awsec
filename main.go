package main

import (
	"log"
	"os"

	"github.com/grimdork/climate/arg"
	"github.com/grimdork/climate/loglines"
)

var (
	version = "dev"
	date    = "undefined"
)

func init() {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
}

const groupCmd = "Commands"

func main() {
	opt := arg.New("awsec", "Store secrets in AWS Parameter Store.")
	opt.SetDefaultHelp(true)
	opt.SetCommand("list", "List keys.", groupCmd, cmdList, []string{"ls"})
	opt.SetCommand("get", "Get value(s) for a key.", groupCmd, cmdGet, []string{"g"})
	opt.SetCommand("set", "Set the value for a key.", groupCmd, cmdSet, []string{"s"})
	opt.SetCommand("tag", "Add tags to a key.", groupCmd, cmdTag, []string{"t"})
	opt.SetCommand("rename", "Rename key.", groupCmd, cmdRename, []string{"ren", "r"})
	opt.SetCommand("remove", "Remove key.", groupCmd, cmdRemove, []string{"rm"})
	opt.SetCommand("backup", "Back up all keys to S3.", groupCmd, cmdBackup, []string{"bak"})
	opt.SetCommand("version", "Show version information.", groupCmd, cmdVersion, []string{"ver", "v"})

	err := opt.Parse(os.Args[1:])
	if err != nil {
		if err == arg.ErrNoArgs {
			return
		}

		if err == arg.ErrRunCommand {
			return
		}

		log.Fatalf("Error parsing arguments: %v\n", err)
	}
}

var pr = loglines.Msg
