package handlogin

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	"git-register-project/internal/models"
	"git-register-project/internal/repository/interface/mocks"
	"git-register-project/internal/servise/login"
)

func TestLoginHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserLogin(ctrl)
	mockMail := mocks.NewMockEmailSender(ctrl)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockToken := mocks.NewMockTokenGenerator(ctrl)

	service := login.NewLoginService(mockRepo, mockMail, mockRedis, mockToken)
	handler := NewLogin(service)

	router := gin.Default()
	router.POST("/login", handler.Login)

	tests := []struct {
		name          string
		body          gin.H
		setupMocks    func()
		expectedCode  int
		expectedError string
	}{
		{
			name: "success login",
			body: gin.H{"email": "ok@example.com", "password": "password123"},
			setupMocks: func() {
				mockRepo.EXPECT().SelectUser(gomock.Any(), "ok@example.com").Return(&models.User{
					ID:       1,
					Email:    "ok@example.com",
					Password: "hashedpassword",
				}, nil)
				mockToken.EXPECT().CheckHash("password123", "hashedpassword").Return(true)
				mockToken.EXPECT().AccessJWT(1).Return("access-token", nil)
				mockToken.EXPECT().RefreshJWT(1).Return("refresh-token", nil)
				mockToken.EXPECT().HashRefreshToken("refresh-token").Return("hashed-refresh")
				mockRepo.EXPECT().InsertRefreshToken(1, "hashed-refresh").Return(nil)
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name: "user not found",
			body: gin.H{"email": "no@example.com", "password": "password123"},
			setupMocks: func() {
				mockRepo.EXPECT().SelectUser(gomock.Any(), "no@example.com").Return(nil, errors.New("not found"))
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "ваш аккаунт не зарегистрирован",
		},
		{
			name: "wrong password",
			body: gin.H{"email": "ok@example.com", "password": "wrongpass"},
			setupMocks: func() {
				mockRepo.EXPECT().SelectUser(gomock.Any(), "ok@example.com").Return(&models.User{
					ID:       1,
					Email:    "ok@example.com",
					Password: "hashedpassword",
				}, nil)
				mockToken.EXPECT().CheckHash("wrongpass", "hashedpassword").Return(false)
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "неверный пароль",
		},
		{
			name: "insert refresh token fails",
			body: gin.H{"email": "ok@example.com", "password": "password123"},
			setupMocks: func() {
				mockRepo.EXPECT().SelectUser(gomock.Any(), "ok@example.com").Return(&models.User{
					ID:       1,
					Email:    "ok@example.com",
					Password: "hashedpassword",
				}, nil)
				mockToken.EXPECT().CheckHash("password123", "hashedpassword").Return(true)
				mockToken.EXPECT().AccessJWT(1).Return("access-token", nil)
				mockToken.EXPECT().RefreshJWT(1).Return("refresh-token", nil)
				mockToken.EXPECT().HashRefreshToken("refresh-token").Return("hashed-refresh")
				mockRepo.EXPECT().InsertRefreshToken(1, "hashed-refresh").Return(errors.New("db error"))
			},
			expectedCode:  http.StatusInternalServerError,
			expectedError: "ошибка вставки refresh токена",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("[%s] expected code %d, got %d", tt.name, tt.expectedCode, w.Code)
			}

			if tt.expectedError != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.expectedError)) {
				t.Errorf("[%s] expected error %q, got %q", tt.name, tt.expectedError, w.Body.String())
			}

			if tt.expectedError == "" && !bytes.Contains(w.Body.Bytes(), []byte("token")) {
				t.Errorf("[%s] expected token in response, got %q", tt.name, w.Body.String())
			}
		})
	}
}
