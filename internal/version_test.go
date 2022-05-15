package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name       string
		beforeFunc func()
		want       string
	}{
		{"no SetVersion", func() {}, "source"},
		{"SetVersion", func() { SetVersion("1.0.0") }, "1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			assert.Equal(t, tt.want, GetVersion())
		})
	}
}
