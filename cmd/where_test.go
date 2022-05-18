package cmd

import (
	"path/filepath"
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/internal"
)

func TestWhere(t *testing.T) {
	srcDir := t.TempDir()

	tests := []struct {
		name          string
		testdata      map[string]interface{}
		beforeFunc    func(*testing.T, string, map[string]interface{})
		want          string
		wantAssertion assert.ComparisonAssertionFunc
		errAssertion  assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			testdata: map[string]interface{}{
				"src_dir": srcDir,
			},
			beforeFunc:   testutil.CreateFile,
			want:         srcDir,
			errAssertion: assert.NoError,
		},
		{
			name: "SrcDirNotFound",
			testdata: map[string]interface{}{
				"src_dir": "",
			},
			beforeFunc:   testutil.CreateFile,
			want:         "",
			errAssertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home := t.TempDir()
			tt.beforeFunc(t, filepath.Join(home, internal.AppName+".json"), tt.testdata)

			_ = internal.InitConfig(internal.WithFile(home))
			var err error
			s := testutil.CaptureOutput(t, func() { err = Where(nil, nil) })
			assert.Equal(t, tt.want, s)
			tt.errAssertion(t, err)
		})
	}
}
