package accessjwt

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		userID    int
		expectErr bool
	}{
		{
			name:      "success",
			secret:    "test-secret",
			userID:    1,
			expectErr: false,
		},
		{
			name:      "midding secret",
			secret:    "",
			userID:    1,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.secret == "" {
				os.Unsetenv("JWT_SECRET")
			} else {
				os.Setenv("JWT_SECRET", tt.secret)
			}
			t.Cleanup(func() {
				os.Unsetenv("JWT_SECRET")
			})
			token, err := AccessJWT((tt.userID))
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			tokens, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(tt.secret), nil
			})
			require.NoError(t, err)
			require.NotNil(t, tokens.Valid)

			claims := tokens.Claims.(jwt.MapClaims)
			require.Equal(t, float64(tt.userID), claims["user_id"])
		})
	}
}
