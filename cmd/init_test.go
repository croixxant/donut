package cmd

import (
	"strings"
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	type args struct {
		srcDir  string
		cfgPath string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{"OK", args{"/home/gopher", ""}, "Configuration file created in", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			s := testutil.CaptureOutput(t, func() { err = Init(tt.args.srcDir, tt.args.cfgPath) })
			assert.True(t, strings.Contains(s, tt.want))
			tt.assertion(t, err)
		})
	}
}
