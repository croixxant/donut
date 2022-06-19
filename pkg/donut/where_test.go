package donut

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDonut_Where(t *testing.T) {
	dir := t.TempDir()

	type fields struct {
		config *Config
	}
	tests := []struct {
		name      string
		fields    fields
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "OK",
			fields: fields{
				config: &Config{
					SrcDir: dir,
				},
			},
			want:      dir + "\n",
			assertion: assert.NoError,
		},
		{
			name: "Error/NoConfig",
			fields: fields{
				config: nil,
			},
			want:      "",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, e := new(bytes.Buffer), new(bytes.Buffer)
			d, _ := New(WithOut(o), WithErr(e))
			d.config = tt.fields.config
			tt.assertion(t, d.Where())
			assert.Equal(t, tt.want, o.String())
		})
	}
}
