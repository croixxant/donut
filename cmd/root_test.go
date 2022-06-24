package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/app"
	"github.com/croixxant/donut/testutil"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "OK"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute()
		})
	}
}

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
			defer app.SetUserHomeDir(home)()

			o, e := new(bytes.Buffer), new(bytes.Buffer)
			cmd := NewRootCmd(o, e)
			args := []string{"init"}
			if tt.opts != nil {
				args = append(args, tt.opts...)
			}
			cmd.SetArgs(args)
			err := cmd.Execute()
			tt.assertion(t, err)
			if err == nil {
				b, err := os.ReadFile(filepath.Join(config, app.Name+".toml"))
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
			want:      filepath.Join(".local", "share", app.Name) + "\n",
			assertion: assert.NoError,
		},
		{
			name:          "OK/WithDefaultConfig",
			opts:          nil,
			defaultConfig: "../testdata/config/basic.toml",
			want:          filepath.Join(".local", "share", app.Name) + "\n",
			assertion:     assert.NoError,
		},
		{
			name:      "OK/WithConfig",
			opts:      []string{"--config", "../testdata/config/basic.toml"},
			want:      filepath.Join(".local", "share", app.Name) + "\n",
			assertion: assert.NoError,
		},
		{
			name:      "Error/WithConfig",
			opts:      []string{"--config", "../testdata/config/dir_not_found.toml"},
			want:      "",
			assertion: assert.Error,
		},
		{
			name:          "Error/Broken",
			opts:          nil,
			defaultConfig: "../testdata/config/broken.toml",
			want:          "",
			assertion:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, config, _ := testutil.CreateBaseDir(t)
			testutil.SetDirEnv(t, home)
			defer app.SetUserHomeDir(home)()

			if tt.defaultConfig != "" {
				defer testutil.CopyFile(t, tt.defaultConfig, filepath.Join(config, app.Name+".toml"))()
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			cmd := NewRootCmd(o, e)
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

func Test_ListCmd(t *testing.T) {
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
			name:      "Error/WithConfig",
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
			cmd := NewRootCmd(o, e)
			args := []string{"list"}
			if tt.opts != nil {
				args = append(args, tt.opts...)
			}
			cmd.SetArgs(args)
			tt.assertion(t, cmd.Execute())
		})
	}
}

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
			cmd := NewRootCmd(o, e)
			args := []string{"apply"}
			if tt.opts != nil {
				args = append(args, tt.opts...)
			}
			cmd.SetArgs(args)
			tt.assertion(t, cmd.Execute())
		})
	}
}
