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
	testutil.CreateFile(t, filepath.Join(dir, mapConfigName+".json"), data)

	tests := []struct {
		name string
		dir  string
		want *MapConfigData
	}{
		{"OK", dir, &MapConfigData{
			Excludes: nil,
			DestDir:  destDir,
			Method:   MethodCopy,
			Maps:     nil,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InitMapConfig(WithFile(mapConfigName, tt.dir))
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
				tt.beforeFunc(t, filepath.Join(dir, mapConfigName+".json"), tt.testdata)
			}
			err := InitMapConfig(WithFile(mapConfigName, dir))
			tt.assertion(t, err)
		})
	}
}

func TestMapConfigData_AbsMaps(t *testing.T) {
	type fields struct {
		Maps map[string]string
	}
	type args struct {
		srcDir  string
		destDir string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "OK",
			fields: fields{
				Maps: map[string]string{
					".config/starship.toml": ".config/starship.toml",
					".zshenv":               "/home/gopher/.zshenv",
				},
			},
			args: args{srcDir: "/home/gopher/.local/share", destDir: "/home/gopher"},
			want: map[string]string{
				"/home/gopher/.local/share/.config/starship.toml": "/home/gopher/.config/starship.toml",
				"/home/gopher/.local/share/.zshenv":               "/home/gopher/.zshenv",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &MapConfigData{
				Maps: tt.fields.Maps,
			}
			assert.Equal(t, tt.want, d.AbsMaps(tt.args.srcDir, tt.args.destDir))
		})
	}
}
