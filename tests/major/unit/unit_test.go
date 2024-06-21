package unit_test

import (
	"Service-oriented-architectures/internal/major"

	"strings"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestLoginCheck(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{"No error", "AwesomeUser123", true},
		{"No error", "User", true},
		{"Just numbers", "12345", true},
		{"Too small", "123", false},
		{"Too big", strings.Repeat("f", 21), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := major.CheckLogin(tt.input)

			require.Equal(t, tt.want, resp)
		})
	}
}

func TestPasswordCheck(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{"No error", "Qwerty12345", true},
		{"No uppercase", "qwerty12345", false},
		{"No lowercase", "QWERTY12345", false},
		{"Small password", "Abc123", false},
		{"Too big password", strings.Repeat("Ab1", 120), false},
		{"Password with space", "Qwerty12345 ", false},
		{"Password with space", "Qwerty12345\n", false},
		{"Password with space", "Qwerty12345\t", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := major.CheckPasswordQuality(tt.input)

			require.Equal(t, tt.want, resp)
		})
	}
}
