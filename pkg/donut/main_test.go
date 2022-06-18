package donut

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "home")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create temp dir")
		os.Exit(1)
	}
	defer os.RemoveAll(home)

	UserHomeDir = home

	os.Exit(m.Run())
}
