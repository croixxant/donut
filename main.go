package main

import (
	"github.com/croixxant/donut/cmd"

	"github.com/croixxant/donut/internal"
)

var version string

func main() {
	if len(version) > 0 {
		internal.SetVersion(version)
	}
	cmd.Execute()
}
