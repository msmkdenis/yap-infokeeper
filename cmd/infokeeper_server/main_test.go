package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Main(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := 2 + 2
			assert.Equal(t, 4, a)
		})
	}
}
