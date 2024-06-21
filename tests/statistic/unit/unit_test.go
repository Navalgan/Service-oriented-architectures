package unit_test

import (
	"Service-oriented-architectures/internal/errors"
	"Service-oriented-architectures/internal/statistic/grpc"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestIsValidUUID(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  bool
	}{
		{"No error", "b84bce46-2f53-11ef-9454-0242ac120002", true},
		{"No error", "31f48f10-07ee-4df6-a026-616c72adc5c8", true},
		{"No error", "e0e02002-5de5-4d9d-8374-b5d8f7d8e6dc", true},
		{"Error", "e0e02002-5de5-4d9d-8374-b5d8f7d8e6dcc", false},
		{"Error", "real uuid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := grpcstatistic.IsValidUUID(tt.input)

			require.Equal(t, tt.want, resp)
		})
	}
}

func TestGetValidOrderBy(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  string
		err   error
	}{
		{"No error", "likes", "Likes", nil},
		{"No error", "Likes", "Likes", nil},
		{"No error", "LiKes", "Likes", nil},
		{"No error", "views", "Views", nil},
		{"No error", "Views", "Views", nil},
		{"No error", "ViewS", "Views", nil},
		{"Error", "l1kes", "", errors.InvalidOrderBy},
		{"Error", "v1ews", "", errors.InvalidOrderBy},
		{"Error", "dislikes", "", errors.InvalidOrderBy},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := grpcstatistic.GetValidOrderBy(tt.input)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, resp)
		})
	}
}
