package app

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()

	type args struct {
		testdata map[string]interface{}
	}
	tests := []struct {
		name      string
		args      args
		want      *Donut
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			args: args{
				testdata: map[string]interface{}{
					"src_dir":  dir,
					"dest_dir": dir,
					"excludes": []string{"example.toml"},
					"method":   "link",
					"mappings": map[string]string{"foo": "bar"},
				},
			},
			want: &Donut{
				config: &Config{
					SrcDir:   dir,
					DestDir:  dir,
					Excludes: []string{"example.toml"},
					Method:   "link",
					Mappings: map[string]string{"foo": "bar"},
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Error",
			args: args{
				testdata: map[string]interface{}{
					"src_dir":  "/home/gopher",
					"dest_dir": "/home/gopher",
				},
			},
			want: &Donut{
				config: &Config{},
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := viper.New()
			if err := v.MergeConfigMap(tt.args.testdata); err != nil {
				t.Fatal(err)
			}
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			tt.want.out, tt.want.err = o, e
			got, err := New(WithConfig(v), WithOut(o), WithErr(e))
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
