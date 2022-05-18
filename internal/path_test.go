package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	tmp := t.TempDir()
	tmpFile := "exists.toml"
	if err := os.WriteFile(filepath.Join(tmp, tmpFile), []byte(""), os.ModePerm); err != nil {
		t.Fatalf("tmpFile not created: %s", err)
	}

	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{"OK", args{s: tmp}, assert.NoError},
		{"IsEmpty", args{s: ""}, assert.Error},
		{"IsFile", args{s: filepath.Join(tmp, tmpFile)}, assert.Error},
		{"IsNotExists", args{s: filepath.Join(tmp, "not_exists")}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, IsDir(tt.args.s))
		})
	}
}

func TestDirOrHome(t *testing.T) {
	tmp, home := t.TempDir(), t.TempDir()

	type args struct {
		dir string
	}
	tests := []struct {
		name       string
		args       args
		beforeFunc func()
		want       string
		assertion  assert.ErrorAssertionFunc
	}{
		{"OK", args{dir: tmp}, func() { t.Setenv("HOME", home) }, tmp, assert.NoError},
		{"DirNotExists/HomeExists", args{dir: filepath.Join(tmp, "not_exists")}, func() { t.Setenv("HOME", home) }, home, assert.NoError},
		{"DirNotExists/HomeNotExists", args{dir: filepath.Join(tmp, "not_exists")}, func() { t.Setenv("HOME", "") }, "", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			got, err := DirOrHome(tt.args.dir)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAbs(t *testing.T) {
	type args struct {
		path    string
		baseDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"OK/Abs", args{"/home/gopher/.local/share", ""}, "/home/gopher/.local/share"},
		{"OK/NotAbs", args{"../.local/share", "/home/gopher/xxx"}, "/home/gopher/.local/share"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Abs(tt.args.path, tt.args.baseDir))
		})
	}
}
