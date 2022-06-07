package cmd

import (
	"path/filepath"
	"testing"

	"github.com/croixxant/donut/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/internal"
)

func TestApply(t *testing.T) {
	home, srcDir := t.TempDir(), t.TempDir()
	type args struct {
		force bool
	}
	tests := []struct {
		name      string
		args      args
		cfgData   map[string]interface{}
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			"OK",
			args{false},
			map[string]interface{}{
				"src_dir":  srcDir,
				"dest_dir": home,
			},
			"",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = internal.InitConfig(internal.WithData(tt.cfgData))
			var err error
			_ = testutil.CaptureOutput(t, func() { err = Apply(tt.args.force) })
			tt.assertion(t, err)
		})
	}
}

func TestLink(t *testing.T) {
	home, srcDir := t.TempDir(), t.TempDir()
	testutil.CreateFile(t, filepath.Join(srcDir, "README.md"), map[string]interface{}{})

	type args struct {
		force bool
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			"OK/NotForce",
			args{false},
			assert.NoError,
		},
		{
			"OK/Force",
			args{true},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := internal.NewMapBuilder(srcDir, home).Build()
			if err != nil {
				t.Fatal(err)
			}

			_ = testutil.CaptureOutput(t, func() { err = Link(m, tt.args.force) })
			tt.assertion(t, err)
		})
	}
}

func TestCopy(t *testing.T) {
	home, srcDir := t.TempDir(), t.TempDir()
	testutil.CreateFile(t, filepath.Join(srcDir, "README.md"), map[string]interface{}{})

	type args struct {
		force bool
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			"OK/NotForce",
			args{false},
			assert.NoError,
		},
		{
			"OK/Force",
			args{true},
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := internal.NewMapBuilder(srcDir, home).Build()
			if err != nil {
				t.Fatal(err)
			}

			_ = testutil.CaptureOutput(t, func() { err = Copy(m, tt.args.force) })
			tt.assertion(t, err)
		})
	}
}
