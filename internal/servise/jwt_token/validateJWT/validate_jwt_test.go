package validatejwt

import (
	generatejwt "git-register-project/internal/servise/jwt_token/accessJWT"
	"os"
	"testing"
)

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		token     string
		expectErr bool
	}{
		{
			name:   "success",
			secret: "test_secret",
			token: func() string {
				os.Setenv("JWT_SECRET", "test_secret")
				token, _ := generatejwt.AccessJWT(1)
				return token
			}(),
			expectErr: false,
		},
		{
			name:      "invalid_secret",
			secret:    "",
			token:     "invalid secret",
			expectErr: true,
		},
		{
			name:      "invalid_token",
			secret:    "test_secret",
			token:     "invalid token",
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

			claims, err := ValidateToken(tt.token)

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if claims["user_id"] == nil {
				t.Fatalf("user_id not found in claims")
			}
		})
	}
}
