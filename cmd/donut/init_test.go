package donut

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/pkg/donut"
	"github.com/croixxant/donut/testutil"
)

func Test_InitCmd(t *testing.T) {
	tests := []struct {
		name      string
		opts      []string
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK",
			opts:      []string{"$HOME/.local/share/donut"},
			want:      "$HOME/.local/share/donut",
			assertion: assert.NoError,
		},
		{
			name:      "OK/NoSpecified",
			opts:      []string{},
			want:      ".local/share/donut",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, config, _ := testutil.CreateBaseDir(t)
			testutil.SetDirEnv(t, home)
			defer donut.SetUserHomeDir(home)()

			o, e := new(bytes.Buffer), new(bytes.Buffer)
			cmd := newRootCmd(o, e)
			args := []string{"init"}
			if tt.opts != nil {
				args = append(args, tt.opts...)
			}
			cmd.SetArgs(args)
			err := cmd.Execute()
			tt.assertion(t, err)
			if err == nil {
				b, err := os.ReadFile(filepath.Join(config, donut.Name+".toml"))
				if err != nil {
					t.Fatal(err)
				}
				data := map[string]string{}
				if err := toml.Unmarshal(b, &data); err != nil {
					t.Fatal(err)
				}
				if tt.want != "" && !filepath.IsAbs(os.ExpandEnv(tt.want)) {
					tt.want = filepath.Join(home, tt.want)
				}
				assert.Equal(t, tt.want, data["src_dir"])
			}
		})
	}
}
