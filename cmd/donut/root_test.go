package donut

import "testing"

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "OK"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute()
		})
	}
}
