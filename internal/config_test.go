package internal

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	home, srcDir := t.TempDir(), t.TempDir()
	data := map[string]interface{}{"src_dir": srcDir}
	cfgName := "config"
	testutil.CreateFile(t, filepath.Join(home, cfgName+".toml"), data)

	tests := []struct {
		name    string
		homeDir string
		want    *ConfigData
	}{
		{
			name:    "OK",
			homeDir: home,
			want:    &ConfigData{SrcDir: srcDir},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InitConfig(WithNameAndPath(cfgName, tt.homeDir))
			assert.Equal(t, tt.want, GetConfig())
		})
	}
}

func TestWriteConfig(t *testing.T) {
	dir := t.TempDir()

	type args struct {
		filename string
		srcDir   string
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"OK",
			args{filepath.Join(dir, "config.toml"), "/home/gopher"},
			"/home/gopher",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InitConfig()
			_, _ = SetConfig("src_dir", tt.args.srcDir)
			tt.assertion(t, WriteConfig(tt.args.filename))
			f, err := os.Open(tt.args.filename)
			if err != nil {
				t.Fatal(err)
			}
			b, err := io.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			assert.True(t, strings.Contains(string(b), tt.want))
		})
	}
}

func TestSetConfig(t *testing.T) {
	srcDir := t.TempDir()

	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      *ConfigData
		assertion assert.ErrorAssertionFunc
	}{
		{"OK", args{"src_dir", srcDir}, &ConfigData{SrcDir: srcDir}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = InitConfig()
			got, err := SetConfig(tt.args.key, tt.args.value)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInitConfig(t *testing.T) {
	srcDir := t.TempDir()

	tests := []struct {
		name       string
		testdata   map[string]interface{}
		beforeFunc func(*testing.T, string, map[string]interface{})
		assertion  assert.ErrorAssertionFunc
	}{
		{
			name:       "OK",
			testdata:   map[string]interface{}{"src_dir": srcDir},
			beforeFunc: testutil.CreateFile,
			assertion:  assert.NoError,
		},
		{
			name:       "FileNotExists",
			testdata:   nil,
			beforeFunc: nil,
			assertion:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			cfgName := "config"
			if tt.beforeFunc != nil {
				tt.beforeFunc(t, filepath.Join(home, cfgName+".toml"), tt.testdata)
			}
			err := InitConfig(WithNameAndPath(cfgName, home))
			tt.assertion(t, err)
		})
	}
}

func TestConfigData_AbsMaps(t *testing.T) {
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
			d := &ConfigData{
				Maps: tt.fields.Maps,
			}
			assert.Equal(t, tt.want, d.AbsMaps(tt.args.srcDir, tt.args.destDir))
		})
	}
}
