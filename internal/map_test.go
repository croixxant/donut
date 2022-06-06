package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_IsSame(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "original.toml"), []byte(""), os.ModePerm); err != nil {
		t.Fatalf("tmpFile not created: %s", err)
	}
	original, err := newFile(filepath.Join(dir, "original.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "original.toml"), filepath.Join(dir, "link.toml")); err != nil {
		t.Fatal(err)
	}
	link, err := newFile(filepath.Join(dir, "link.toml"))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		path string
	}
	tests := []struct {
		name      string
		file      *File
		args      args
		want      bool
		assertion assert.ErrorAssertionFunc
	}{
		{"IsSymlink/Same", link, args{filepath.Join(dir, "original.toml")}, true, assert.NoError},
		{"IsNotSymlink/Same", original, args{filepath.Join(dir, "original.toml")}, true, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.file
			got, err := f.IsSame(tt.args.path)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
