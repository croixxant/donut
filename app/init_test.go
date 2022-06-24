package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/croixxant/donut/testutil"
)

func TestDonut_Init(t *testing.T) {
	testutil.SetDirEnv(t, UserHomeDir)

	type fields struct {
		config *Config
	}
	type args struct {
		srcDir  string
		cfgPath string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		defaultConfig string
		want          map[string]interface{}
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: "",
			},
			want: map[string]interface{}{
				"src_dir":  filepath.Join(UserHomeDir, ".local", "share", Name),
				"dest_dir": UserHomeDir,
				"method":   MethodCopy,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/SrcDirSpecified",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  filepath.Join(UserHomeDir, Name),
				cfgPath: "",
			},
			want: map[string]interface{}{
				"src_dir":  filepath.Join(UserHomeDir, Name),
				"dest_dir": UserHomeDir,
				"method":   MethodCopy,
			},
			assertion: assert.NoError,
		},
		{
			name: "OK/CfgPathSpecified",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: filepath.Join(UserHomeDir, ".config", Name+".toml"),
			},
			want: map[string]interface{}{
				"src_dir":  filepath.Join(UserHomeDir, ".local", "share", Name),
				"dest_dir": UserHomeDir,
				"method":   MethodCopy,
			},
			assertion: assert.NoError,
		},
		{
			name: "Error/NoNeed",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: "",
			},
			defaultConfig: "../testdata/config/basic.toml",
			want:          map[string]interface{}{},
			assertion:     assert.Error,
		},
		{
			name: "Error/FileBroken",
			fields: fields{
				config: nil,
			},
			args: args{
				srcDir:  "",
				cfgPath: "",
			},
			defaultConfig: "../testdata/config/broken.toml",
			want:          map[string]interface{}{},
			assertion:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.defaultConfig != "" {
				defer testutil.CopyFile(t, tt.defaultConfig, filepath.Join(UserHomeDir, ".config", Name+".toml"))()
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			cfgPath := filepath.Join(UserHomeDir, ".config", Name, Name+".toml")
			if tt.args.cfgPath != "" {
				cfgPath = tt.args.cfgPath
			}
			err := d.Init(tt.args.srcDir, tt.args.cfgPath)
			defer func() { os.Remove(cfgPath) }()
			tt.assertion(t, err)
			if err != nil {
				return
			}
			if diff := cmp.Diff(tt.want, testutil.ReadToml(t, cfgPath)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
