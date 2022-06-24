package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/croixxant/donut/app"
)

func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "home")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create temp dir")
		os.Exit(1)
	}
	defer os.RemoveAll(home)

	app.UserHomeDir = home

	os.Exit(m.Run())
}
