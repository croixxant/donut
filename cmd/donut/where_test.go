package donut

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/pkg/donut"
	"github.com/croixxant/donut/testutil"
)

func Test_WhereCmd(t *testing.T) {
	tests := []struct {
		name          string
		opts          []string
		defaultConfig string
		want          string
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name:      "OK/WithNoConfig",
			opts:      nil,
			want:      filepath.Join(".local", "share", donut.Name),
			assertion: assert.NoError,
		},
		{
			name:          "OK/WithDefaultConfig",
			opts:          nil,
			defaultConfig: "../../testdata/config/basic.toml",
			want:          filepath.Join(".local", "share", donut.Name),
			assertion:     assert.NoError,
		},
		{
			name:      "OK/WithConfig",
			opts:      []string{"--config", "../../testdata/config/basic.toml"},
			want:      filepath.Join(".local", "share", donut.Name),
			assertion: assert.NoError,
		},
		{
			name:      "Error/WithConfig",
			opts:      []string{"--config", "../../testdata/config/dir_not_found.toml"},
			want:      "",
			assertion: assert.Error,
		},
		{
			name:          "Error/Broken",
			opts:          nil,
			defaultConfig: "../../testdata/config/broken.toml",
			want:          "",
			assertion:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, config, _ := testutil.CreateBaseDir(t)
			testutil.SetDirEnv(t, home)
			defer donut.SetUserHomeDir(home)()

			if tt.defaultConfig != "" {
				defer testutil.CopyFile(t, tt.defaultConfig, filepath.Join(config, donut.Name+".toml"))()
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			cmd := newRootCmd(o, e)
			args := []string{"where"}
			if tt.opts != nil {
				args = append(args, tt.opts...)
			}
			cmd.SetArgs(args)
			tt.assertion(t, cmd.Execute())
			if !filepath.IsAbs(tt.want) && tt.want != "" {
				tt.want = filepath.Join(home, tt.want)
			}
			assert.Equal(t, tt.want, o.String())
		})
	}
}
