package checkhash

import (
	hashrefreshtoken "git-register-project/internal/servise/jwt_token/hashRefreshToken"
	"testing"
)

func TestCheckRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		tokenHash string
		expected  bool
	}{
		{
			name:      "success",
			token:     "test_token",
			tokenHash: hashrefreshtoken.HashRefreshToken("test_token"),
			expected:  true,
		},
		{
			name:      "invalid token",
			token:     "",
			tokenHash: hashrefreshtoken.HashRefreshToken("test_token"),
			expected:  false,
		},
		{
			name:      "invalid hash",
			token:     "test_token",
			tokenHash: "",
			expected:  false,
		},
		{
			name:      "token mismatch",
			token:     "other_token",
			tokenHash: hashrefreshtoken.HashRefreshToken("test_token"),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckRefreshToken(tt.token, tt.tokenHash); got != tt.expected {
				t.Errorf("CheckRefreshToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}
