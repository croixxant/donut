package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
		want string
	}{
		{"NoSetVersion", func() {}, "source"},
		{"SetVersion", func() { SetVersion("1.0.0") }, "1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			assert.Equal(t, tt.want, GetVersion())
		})
	}
}
