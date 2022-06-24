package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/app"
	"github.com/croixxant/donut/testutil"
)

func Test_ApplyCmd(t *testing.T) {
	tests := []struct {
		name      string
		opts      []string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK/WithNoConfig",
			opts:      nil,
			assertion: assert.NoError,
		},
		{
			name:      "Error/WithConfig/DirNotFound",
			opts:      []string{"--config", "../testdata/config/dir_not_found.toml"},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, _, _ := testutil.CreateBaseDir(t)
			testutil.SetDirEnv(t, home)
			defer app.SetUserHomeDir(home)()

			o, e := new(bytes.Buffer), new(bytes.Buffer)
			cmd := newRootCmd(o, e)
			args := []string{"apply"}
			if tt.opts != nil {
				args = append(args, tt.opts...)
			}
			cmd.SetArgs(args)
			tt.assertion(t, cmd.Execute())
		})
	}
}
