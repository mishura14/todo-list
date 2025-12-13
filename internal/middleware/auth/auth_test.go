package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	generatejwt "git-register-project/internal/servise/jwt_token/accessJWT"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test_secret")
	defer os.Unsetenv("JWT_SECRET")

	validToken, err := generatejwt.AccessJWT(123)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name       string
		authHeader string
		statusCode int
		expectedID interface{}
	}{
		{
			name:       "valid token",
			authHeader: "Bearer " + validToken,
			statusCode: http.StatusOK,
			expectedID: float64(123),
		},
		{
			name:       "no header",
			authHeader: "",
			statusCode: http.StatusUnauthorized,
			expectedID: nil,
		},
		{
			name:       "invalid token",
			authHeader: "Bearer invalidtoken",
			statusCode: http.StatusUnauthorized,
			expectedID: nil,
		},
		{
			name:       "wrong prefix",
			authHeader: "Token " + validToken,
			statusCode: http.StatusUnauthorized,
			expectedID: nil,
		},
	}

	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusOK, gin.H{"user_id": nil})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.statusCode, w.Code)
			}

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			if err != nil {
				t.Fatalf("%s: failed to parse response: %v", tt.name, err)
			}

			userID, ok := resp["user_id"]
			if tt.expectedID == nil {
				if ok && userID != nil {
					t.Errorf("%s: expected no user_id, got %v", tt.name, userID)
				}
			} else {
				if userID != tt.expectedID {
					t.Errorf("%s: expected user_id %v, got %v", tt.name, tt.expectedID, userID)
				}
			}
		})
	}
}
