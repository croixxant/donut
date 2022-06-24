package main

import "github.com/croixxant/donut/cmd"

var version string

func main() {
	if len(version) > 0 {
		cmd.SetVersion(version)
	}
	cmd.Execute()
}
