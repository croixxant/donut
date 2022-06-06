package cmd

import (
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/internal"
)

func TestWhere(t *testing.T) {
	srcDir := t.TempDir()

	tests := []struct {
		name      string
		testdata  map[string]interface{}
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "OK",
			testdata:  map[string]interface{}{"src_dir": srcDir},
			want:      srcDir,
			assertion: assert.NoError,
		},
		{
			name:      "IsNotDir",
			testdata:  map[string]interface{}{},
			want:      "",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = internal.InitConfig(internal.WithData(tt.testdata))
			var err error
			s := testutil.CaptureOutput(t, func() { err = Where() })
			assert.Equal(t, tt.want, s)
			tt.assertion(t, err)
		})
	}
}
