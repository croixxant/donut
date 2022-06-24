package main

import (
	"os"

	"github.com/pterm/pterm"

	"github.com/croixxant/donut/cmd"
)

var version string

func main() {
	if len(version) > 0 {
		cmd.SetVersion(version)
	}
	if err := cmd.NewRootCmd(os.Stdout, os.Stderr).Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
