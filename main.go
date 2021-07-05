package main

import (
	"fmt"
	"os"

	"github.com/Urethramancer/signor/opt"
)

var o struct {
	opt.DefaultHelp
	List   ListCmd   `command:"list" help:"List keys." aliases:"ls"`
	Get    GetCmd    `command:"get" help:"Get value(s) for a key." aliases:"g"`
	Set    SetCmd    `command:"set" help:"Set the value for a key." aliases:"s"`
	Tag    TagCmd    `command:"tag" help:"Add tags to a key." aliases:"t"`
	Rename RenameCmd `command:"rename" help:"Rename key." aliases:"ren,r"`
	Remove RemoveCmd `command:"remove" help:"Remove key." aliases:"rm"`
}

func init() {
	os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
}

func main() {
	a := opt.Parse(&o)
	if o.Help {
		a.Usage()
		return
	}

	err := a.RunCommand(false)
	if err != nil {
		if err == opt.ErrNoCommand {
			a.Usage()
			return
		}

		pr("Error running command: %s", err.Error())
		os.Exit(2)
	}
}

func pr(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}
