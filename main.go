package main

import "github.com/croixxant/donut/cmd/donut"

var version string

func main() {
	if len(version) > 0 {
		donut.SetVersion(version)
	}
	donut.Execute()
}
