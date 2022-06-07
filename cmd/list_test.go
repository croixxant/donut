package cmd

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"

	"github.com/croixxant/donut/internal"
)

func TestList(t *testing.T) {
	home, srcDir := t.TempDir(), t.TempDir()
	if err := os.MkdirAll(filepath.Join(srcDir, ".config"), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	// create files to srcDir
	testutil.CreateFile(t, filepath.Join(srcDir, "README.md"), map[string]interface{}{})
	testutil.CreateFile(t, filepath.Join(srcDir, "foo.toml"), map[string]interface{}{})
	testutil.CreateFile(t, filepath.Join(srcDir, ".zshenv"), map[string]interface{}{})
	testutil.CreateFile(t, filepath.Join(srcDir, ".config/foo.toml"), map[string]interface{}{})
	// symlink README.md -> foo.toml
	if err := os.Symlink(filepath.Join(srcDir, "README.md"), filepath.Join(home, "foo.toml")); err != nil {
		t.Fatal(err)
	}
	// symlink .zshenv -> .zshenv
	if err := os.Symlink(filepath.Join(srcDir, ".zshenv"), filepath.Join(home, ".zshenv")); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		cfgData   map[string]interface{}
		want      [][]string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"OK",
			map[string]interface{}{
				"src_dir":  srcDir,
				"dest_dir": home,
				"method":   "link",
				"maps": map[string]string{
					"foo.toml": "changed.toml",
				},
			},
			[][]string{
				{"", filepath.Join(srcDir, ".config/foo.toml"), filepath.Join(home, ".config/foo.toml")},
				{"✔", filepath.Join(srcDir, ".zshenv"), filepath.Join(home, ".zshenv")},
				{"", filepath.Join(srcDir, "README.md"), filepath.Join(home, "README.md")},
				{"", filepath.Join(srcDir, "foo.toml"), filepath.Join(home, "changed.toml")},
			},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = internal.InitConfig(internal.WithData(tt.cfgData))
			var err error
			s := testutil.CaptureOutput(t, func() { err = List() })
			b := bufio.NewReader(bytes.NewBufferString(s))
			for i := 0; i < 2; i++ { // skip box line
				if _, _, err := b.ReadLine(); err != nil {
					t.Fatal(err)
				}
			}

			for _, items := range tt.want {
				line, _, err := b.ReadLine()
				if err == io.EOF {
					break
				} else if err != nil {
					t.Fatal(err)
				}
				for _, item := range items {
					assert.Contains(t, string(line), item)
				}
			}

			tt.assertion(t, err)
		})
	}
}
