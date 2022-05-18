package internal

import (
	"reflect"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestExpandEnvFunc(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	tests := []struct {
		name string
		f, t reflect.Value
		want interface{}
		err  error
	}{
		{"OK", reflect.ValueOf("$HOME/.local/share"), reflect.ValueOf(""),
			tmp + "/.local/share", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapstructure.DecodeHookExec(ExpandEnvFunc(), tt.f, tt.t)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}