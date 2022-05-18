package internal

import (
	"path/filepath"
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetMapConfig(t *testing.T) {
	dir, destDir := t.TempDir(), t.TempDir()
	data := map[string]interface{}{"dest_dir": destDir}
	testutil.CreateFile(t, filepath.Join(dir, MapConfigName+".json"), data)

	tests := []struct {
		name string
		dir  string
		want *MapConfigData
	}{
		{"OK", dir, &MapConfigData{
			Excludes: nil,
			DestDir:  destDir,
			Files:    nil,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InitMapConfig(WithFile(tt.dir))
			assert.Equal(t, tt.want, GetMapConfig())
		})
	}
}

func TestInitMapConfig(t *testing.T) {
	destDir := t.TempDir()

	tests := []struct {
		name       string
		testdata   map[string]interface{}
		beforeFunc func(*testing.T, string, map[string]interface{})
		assertion  assert.ErrorAssertionFunc
	}{
		{
			name:       "OK",
			testdata:   map[string]interface{}{"dest_dir": destDir},
			beforeFunc: testutil.CreateFile,
			assertion:  assert.NoError,
		},
		{
			name:       "FileNotExists/HomeExists",
			testdata:   nil,
			beforeFunc: nil,
			assertion:  assert.NoError,
		},
		{
			name:     "FileNotExists/HomeNotExists",
			testdata: nil,
			beforeFunc: func(*testing.T, string, map[string]interface{}) {
				t.Setenv("HOME", "")
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.beforeFunc != nil {
				tt.beforeFunc(t, filepath.Join(dir, MapConfigName+".json"), tt.testdata)
			}
			err := InitMapConfig(WithFile(dir))
			tt.assertion(t, err)
		})
	}
}

func TestMapConfigData_AbsFiles(t *testing.T) {
	type fields struct {
		Files []Map
	}
	type args struct {
		srcDir  string
		destDir string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Map
	}{
		{
			name: "OK",
			fields: fields{
				Files: []Map{
					{Src: ".config/starship.toml", Dest: ".config/starship.toml"},
					{Src: ".zshenv", Dest: "/home/gopher/.zshenv"},
				},
			},
			args: args{srcDir: "/home/gopher/.local/share", destDir: "/home/gopher"},
			want: []Map{
				{Src: "/home/gopher/.local/share/.config/starship.toml", Dest: "/home/gopher/.config/starship.toml"},
				{Src: "/home/gopher/.local/share/.zshenv", Dest: "/home/gopher/.zshenv"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MapConfigData{
				Files: tt.fields.Files,
			}
			assert.Equal(t, tt.want, d.AbsFiles(tt.args.srcDir, tt.args.destDir))
		})
	}
}
