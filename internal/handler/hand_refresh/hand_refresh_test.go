package handrefresh_test

import (
	"bytes"
	"encoding/json"
	handrefresh "git-register-project/internal/handler/hand_refresh"
	"git-register-project/internal/models"
	"git-register-project/internal/repository/interface/mocks"
	refreshtoken "git-register-project/internal/servise/refreshToken"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRefreshHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       models.Refresh
		setupMocks func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken, redis *mocks.MockRedisClient)
		wantStatus int
	}{
		{
			name: "success",
			body: models.Refresh{Token: "valid refresh"},

			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken, redis *mocks.MockRedisClient) {
				id := 1
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(id)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), id).Return("hash-refresh", nil)
				tg.EXPECT().CheckHash("valid refresh", "hash-refresh").Return(true)
				tg.EXPECT().AccessJWT(id).Return("new-Access", nil)
				tg.EXPECT().RefreshJWT(id).Return("new-Refresh", nil)
				tg.EXPECT().HashRefreshToken("new-Refresh").Return("new-hash")
				repo.EXPECT().UpdateRefreshToken(gomock.Any(), id, "new-hash", gomock.AssignableToTypeOf(time.Time{})).Return(nil)
			},

			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tg := mocks.NewMockTokenGenerator(ctrl)
			repo := mocks.NewMockRefreshtoken(ctrl)
			redis := mocks.NewMockRedisClient(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(tg, repo, redis)
			}

			service := refreshtoken.NewRefreshService(tg, redis, repo)
			handler := handrefresh.NewHandlerRefresh(service)

			router := gin.Default()
			router.POST("/refresh", handler.Refresh)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
