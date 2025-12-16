package hashrefreshtoken

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

// Тестируем HashRefreshToken через табличный тест
func TestHashRefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "simple token",
			token:    "mytoken123",
			expected: func() string { h := sha256.Sum256([]byte("mytoken123")); return hex.EncodeToString(h[:]) }(),
		},
		{
			name:     "empty token",
			token:    "",
			expected: func() string { h := sha256.Sum256([]byte("")); return hex.EncodeToString(h[:]) }(),
		},
		{
			name:  "long token",
			token: "this_is_a_very_long_token_which_exceeds_bcrypt_limit_but_sha256_can_handle_it_without_problems",
			expected: func() string {
				h := sha256.Sum256([]byte("this_is_a_very_long_token_which_exceeds_bcrypt_limit_but_sha256_can_handle_it_without_problems"))
				return hex.EncodeToString(h[:])
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashRefreshToken(tt.token)
			if got != tt.expected {
				t.Errorf("HashRefreshToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}
