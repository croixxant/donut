package app

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"
)

func TestNewConfig(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "example.toml")
	testutil.WriteToml(t, file, map[string]interface{}{
		"src_dir":  dir,
		"dest_dir": dir,
	})
	type args struct {
		opts []ConfigOption
	}
	type want struct {
		SrcDir  string
		DestDir string
	}
	tests := []struct {
		name      string
		args      args
		want      want
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK/WithDefault",
			args: args{
				opts: []ConfigOption{WithDefault()},
			},
			want: want{
				SrcDir:  filepath.Join(UserHomeDir, ".local", "share", Name),
				DestDir: UserHomeDir,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/WithData",
			args: args{
				opts: []ConfigOption{WithData(map[string]interface{}{
					"src_dir":  dir,
					"dest_dir": dir,
				})},
			},
			want: want{
				SrcDir:  dir,
				DestDir: dir,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/WithNameAndPath",
			args: args{
				opts: []ConfigOption{WithNameAndPath("example", dir)},
			},
			want: want{
				SrcDir:  dir,
				DestDir: dir,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/WithFile",
			args: args{
				opts: []ConfigOption{WithFile(file)},
			},
			want: want{
				SrcDir:  dir,
				DestDir: dir,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/WithNameAndPath",
			args: args{
				opts: []ConfigOption{WithNameAndPath("example", "/home/gopher")},
			},
			want:      want{},
			assertion: assert.Error,
		},
		{
			name: "Error/WithFile",
			args: args{
				opts: []ConfigOption{WithFile("")},
			},
			want:      want{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.opts...)
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want.SrcDir, got.GetString("src_dir"))
				assert.Equal(t, tt.want.DestDir, got.GetString("dest_dir"))
			}
		})
	}
}

func TestConfig_absMappings(t *testing.T) {
	dir := t.TempDir()

	type fields struct {
		SrcDir   string
		DestDir  string
		Mappings map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "OK",
			fields: fields{
				SrcDir:  dir,
				DestDir: dir,
				Mappings: map[string]string{
					"example.toml": "example.toml",
					".example":     "/home/gopher/.example",
				},
			},
			want: map[string]string{
				filepath.Join(dir, "example.toml"): filepath.Join(dir, "example.toml"),
				filepath.Join(dir, ".example"):     "/home/gopher/.example",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Config{
				SrcDir:   tt.fields.SrcDir,
				DestDir:  tt.fields.DestDir,
				Mappings: tt.fields.Mappings,
			}
			assert.Equal(t, tt.want, d.absMappings())
		})
	}
}
