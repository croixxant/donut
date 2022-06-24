package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"
)

func Test_isDir(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "example")
	if _, err := os.Create(file); err != nil {
		t.Fatal(err)
	}

	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				s: dir,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/IsEmpty",
			args: args{
				s: "",
			},
			assertion: assert.Error,
		},
		{
			name: "Error/IsNotExist",
			args: args{
				s: "/home/gopher/example",
			},
			assertion: assert.Error,
		},
		{
			name: "Error/IsFile",
			args: args{
				s: file,
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, isDir(tt.args.s))
		})
	}
}

func Test_copyFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "example.toml")
	testutil.WriteToml(t, file, map[string]interface{}{"foo": "bar"})

	type args struct {
		src  string
		dest string
	}
	tests := []struct {
		name      string
		args      args
		want      map[string]interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				src:  file,
				dest: filepath.Join(dir, "example_copy.toml"),
			},
			want:      map[string]interface{}{"foo": "bar"},
			assertion: assert.NoError,
		},
		{
			name: "Error/CannotOpen",
			args: args{
				src:  "/home/gopher/example",
				dest: "/home/gopher/example",
			},
			want:      nil,
			assertion: assert.Error,
		},
		{
			name: "Error/CannotCreate",
			args: args{
				src:  file,
				dest: "/home/gopher/example",
			},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := copyFile(tt.args.src, tt.args.dest)
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, testutil.ReadToml(t, tt.args.dest))
			}
		})
	}
}
