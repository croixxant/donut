package donut

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"
)

func TestDonut_List(t *testing.T) {
	srcFiles := map[string]string{
		"../../testdata/data/.example":     ".example",
		"../../testdata/data/example.toml": "example.toml",
		"../../testdata/data/example.yaml": ".config/example.yaml",
		"../../testdata/data/example.json": "example.json",
		"../../testdata/data/example.ini":  "example.ini",
		"../../testdata/data/.zshrc":       ".zshrc",
		"../../testdata/data/.zprofile":    ".zprofile",
		"../../testdata/data/.zshenv":      ".zshenv",
		"../../testdata/data/.zlogin":      ".zlogin",
	}
	destFiles := map[string]string{"../../testdata/data/.example": ".example"}
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
	tests := []struct {
		name      string
		srcFiles  map[string]string // map[testdata]srcdata
		destFiles map[string]string // map[testdata]destdata
		links     map[string]string // map[srcdata]destdata
		fields    fields
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
			assertion: assert.NoError,
		},
		{
			name: "Error/NoConfig",
			fields: fields{
				config: nil,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, _, data := testutil.CreateBaseDir(t)
			testutil.CreateDirs(t, filepath.Join(data, ".config"))
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
			tt.assertion(t, d.List())
			// TODO: assert output
		})
	}
}
