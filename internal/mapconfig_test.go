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
