package internal

import (
	"path/filepath"
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	home, srcDir := t.TempDir(), t.TempDir()
	data := map[string]interface{}{"src_dir": srcDir}
	testutil.CreateFile(t, filepath.Join(home, AppName+".json"), data)

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
			_ = InitConfig(WithFile(tt.homeDir))
			assert.Equal(t, tt.want, GetConfig())
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
			if tt.beforeFunc != nil {
				tt.beforeFunc(t, filepath.Join(home, AppName+".json"), tt.testdata)
			}
			err := InitConfig(WithFile(home))
			tt.assertion(t, err)
		})
	}
}
