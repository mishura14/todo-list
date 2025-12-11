package handler_register_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	handler_register "git-register-project/internal/handler/hand_register/register"
	"git-register-project/internal/repository/mocks"
	"git-register-project/internal/servise/register/register"
)

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       gin.H
		setupMocks func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient)
		status     int
		message    string
	}{
		{
			name: "success",
			body: gin.H{"email": "ok@example.com", "password": "mishura_14", "name": "mishura"},

			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				mail.EXPECT().SendConfremRegister("ok@example.com", gomock.Any()).Return(nil)
				redis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},

			status:  http.StatusOK,
			message: "код регистрации отправлен",
		},

		{
			name: "invalid email format",
			body: gin.H{"email": "bad email", "password": "mishura_14", "name": "mishura"},

			setupMocks: nil,

			status:  http.StatusBadRequest,
			message: "неверный формат email",
		},
		{
			name: "email exists",
			body: gin.H{"email": "ok@example.com", "password": "mishura_14", "name": "mishura"},

			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(true, nil)
			},

			status:  http.StatusConflict,
			message: "email уже зарегистрирован",
		},
		{
			name: "valid password",
			body: gin.H{"email": "ok@example.com", "password": "mis14", "name": "mishura"},

			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
			},

			status:  http.StatusBadRequest,
			message: "пароль не соответствует требованиям",
		},
		{
			name: "invalid mail",
			body: gin.H{"email": "ok@example.com", "password": "mishura_14", "name": "mishura"},

			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				mail.EXPECT().SendConfremRegister("ok@example.com", gomock.Any()).Return(errors.New("smtp error"))
			},

			status:  http.StatusInternalServerError,
			message: "не удалось отправить письмо",
		},
		{
			name: "invalid redis",
			body: gin.H{"email": "ok@example.com", "password": "mishura_14", "name": "mishura"},

			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				mail.EXPECT().SendConfremRegister("ok@example.com", gomock.Any()).Return(nil)
				redis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
			},

			status:  http.StatusInternalServerError,
			message: "ошибка сохранения в redis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUserRegister(ctrl)
			mail := mocks.NewMockEmailSender(ctrl)
			redis := mocks.NewMockRedisClient(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(repo, mail, redis)
			}

			service := register.NewRegisterService(repo, mail, redis)
			handler := handler_register.NewRegister(service)

			router := gin.Default()
			router.POST("/register", handler.Register)

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, w.Code)
			}

			if !bytes.Contains(w.Body.Bytes(), []byte(tt.message)) {
				t.Errorf("expected message %q, got %q", tt.message, w.Body.String())
			}
		})
	}
}
