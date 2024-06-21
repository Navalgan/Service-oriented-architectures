package unit_test

import (
	"Service-oriented-architectures/internal/common"
	"github.com/stretchr/testify/require"

	"time"

	"testing"
)

func TestJwtTokenWork(t *testing.T) {
	var tests = []struct {
		name        string
		jwtToken    []byte
		inputUserID string
		inputLogin  string
		wantUserID  string
		wantLogin   string
	}{
		{"No error", []byte("token"), "user_id123", "AwesomeUser", "user_id123", "AwesomeUser"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signedToken, err := common.NewToken(tt.jwtToken, tt.inputUserID, tt.inputLogin, time.Second*100)

			require.Equal(t, nil, err)

			unsignedToken, err := common.VerifyToken(tt.jwtToken, signedToken)

			require.Equal(t, nil, err)
			require.Equal(t, tt.wantUserID, unsignedToken.UserID)
			require.Equal(t, tt.wantLogin, unsignedToken.Login)
		})
	}
}

func TestJwtTokenError(t *testing.T) {
	var tests = []struct {
		name        string
		jwtToken    []byte
		inputUserID string
		inputLogin  string
	}{
		{"Error", []byte("token"), "user_id123", "AwesomeUser"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := "User's ID"
			userLogin := "User's Login"

			signedToken, err := common.NewToken(tt.jwtToken, userID, userLogin, time.Second*100)

			require.Equal(t, nil, err)

			unsignedToken, err := common.VerifyToken(tt.jwtToken, signedToken)

			require.Equal(t, nil, err)
			require.NotEqual(t, tt.inputUserID, unsignedToken.UserID)
			require.NotEqual(t, tt.inputLogin, unsignedToken.Login)
		})
	}
}
