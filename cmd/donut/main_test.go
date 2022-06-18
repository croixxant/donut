package donut

import (
	"fmt"
	"os"
	"testing"

	"github.com/croixxant/donut/pkg/donut"
)

func TestMain(m *testing.M) {
	home, err := os.MkdirTemp("", "home")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create temp dir")
		os.Exit(1)
	}
	defer os.RemoveAll(home)

	donut.UserHomeDir = home

	os.Exit(m.Run())
}
