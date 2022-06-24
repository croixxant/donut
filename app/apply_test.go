package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"
)

func TestDonut_Apply(t *testing.T) {
	srcFiles := map[string]string{
		"../testdata/data/.example":     ".example",
		"../testdata/data/example.toml": "example.toml",
		"../testdata/data/example.yaml": ".config/example.yaml",
		"../testdata/data/example.json": "example.json",
		"../testdata/data/example.ini":  "example.ini",
		"../testdata/data/.zshrc":       ".zshrc",
		"../testdata/data/.zprofile":    ".zprofile",
		"../testdata/data/.zshenv":      ".zshenv",
		"../testdata/data/.zlogin":      ".zlogin",
		"../testdata/data/.empty":       ".git/.empty",
		"../testdata/data/.gitconfig":   ".gitconfig",
	}
	destFiles := map[string]string{"../testdata/data/.example": ".example"}
	links := map[string]string{
		".zprofile": ".zprofile",
		".zshrc":    ".zlogin",
	}
	mappings := map[string]string{
		"example.toml": ".config/example.toml",
	}
	excludes := []string{"example.json"}

	type fields struct {
		config *Config
	}
	type args struct {
		force bool
	}
	tests := []struct {
		name      string
		srcFiles  map[string]string // map[testdata]srcdata
		destFiles map[string]string // map[testdata]destdata
		links     map[string]string // map[srcdata]destdata
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK/Copy",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodCopy,
					Mappings: mappings,
				},
			},
			args: args{
				force: false,
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Link",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: mappings,
				},
			},
			args: args{
				force: false,
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Copy/Force",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodCopy,
					Mappings: mappings,
				},
			},
			args: args{
				force: true,
			},
			assertion: assert.NoError,
		},
		{
			name:      "OK/Link/Force",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: mappings,
				},
			},
			args: args{
				force: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/NoConfig",
			fields: fields{
				config: nil,
			},
			args: args{
				force: false,
			},
			assertion: assert.Error,
		},
		{
			name:      "Error/Link/OtherLinkExists",
			srcFiles:  srcFiles,
			destFiles: destFiles,
			links:     links,
			fields: fields{
				config: &Config{
					Excludes: excludes,
					Method:   MethodLink,
					Mappings: map[string]string{
						"example.toml": ".config/example.toml",
						".zshenv":      ".zshrc",
					},
				},
			},
			args: args{
				force: true,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, _, data := testutil.CreateBaseDir(t)
			testutil.CreateDirs(t, filepath.Join(data, ".config"), filepath.Join(data, ".git"))
			for s, d := range tt.srcFiles {
				defer testutil.CopyFile(t, s, filepath.Join(data, d))()
			}
			for s, d := range tt.destFiles {
				defer testutil.CopyFile(t, s, filepath.Join(home, d))()
			}
			for s, d := range tt.links {
				if err := os.Symlink(filepath.Join(data, s), filepath.Join(home, d)); err != nil {
					t.Fatal(err)
				}
				defer os.Remove(d)
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			if tt.fields.config != nil {
				d.config = tt.fields.config
				d.config.SrcDir = data
				d.config.DestDir = home
			}
			tt.assertion(t, d.Apply(tt.args.force))
			// TODO: assert output
		})
	}
}
