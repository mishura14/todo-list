package handler_comfirm_register_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	handler_comfirm_register "git-register-project/internal/handler/hand_register/confirm_register"
	"git-register-project/internal/models"
	"git-register-project/internal/repository/interface/mocks"
	comfirm_register "git-register-project/internal/servise/register/confirm_register"
)

func TestRegisterConfirm(t *testing.T) {

	tests := []struct {
		name       string
		code       string
		setupMocks func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient)
		status     int
		message    string
	}{
		{
			name: "success",
			code: "1234",
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				userJSON, _ := json.Marshal(models.UserRedis{
					Name:     "mishura",
					Email:    "ok@example.com",
					Password: "hashed_pass",
				})

				redis.EXPECT().Get(gomock.Any(), "1234").Return([]byte(userJSON), nil)
				repo.EXPECT().CreateUser(gomock.Any()).Return(nil)
				redis.EXPECT().Del(gomock.Any(), "1234").Return(nil)
			},
			status:  http.StatusOK,
			message: "регистрация успешно завершена",
		},

		{
			name: "timeout code",
			code: "1234",
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				redis.EXPECT().Get(gomock.Any(), "1234").Return(nil, errors.New("timeout"))
			},
			status:  http.StatusBadRequest,
			message: "код подтверждения устарел",
		},

		{
			name: "bad json",
			code: "1234",
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				redis.EXPECT().Get(gomock.Any(), "1234").Return([]byte("not-json"), nil)
			},
			status:  http.StatusInternalServerError,
			message: "ошибка обработки JSON",
		},

		{
			name: "create user error",
			code: "1234",
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				userJSON, _ := json.Marshal(models.UserRedis{
					Name:     "mishura",
					Email:    "ok@example.com",
					Password: "hashed_pass",
				})

				redis.EXPECT().Get(gomock.Any(), "1234").Return([]byte(userJSON), nil)
				repo.EXPECT().CreateUser(gomock.Any()).Return(errors.New("create error"))
			},
			status:  http.StatusInternalServerError,
			message: "ошибка создания пользователя",
		},

		{
			name: "delete code error",
			code: "1234",
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				userJSON, _ := json.Marshal(models.UserRedis{
					Name:     "mishura",
					Email:    "ok@example.com",
					Password: "hashed_pass",
				})

				redis.EXPECT().Get(gomock.Any(), "1234").Return([]byte(userJSON), nil)
				repo.EXPECT().CreateUser(gomock.Any()).Return(nil)
				redis.EXPECT().Del(gomock.Any(), "1234").Return(errors.New("del error"))
			},
			status:  http.StatusInternalServerError,
			message: "ошибка удаления кода подтверждения",
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

			service := comfirm_register.NewConfirmRegisterService(repo, mail, redis)
			handler := handler_comfirm_register.NewConfirmRegister(service)

			router := gin.Default()
			router.POST("/confirm", handler.Confirm_register)

			bodyJSON, _ := json.Marshal(map[string]string{"code": tt.code})
			req, _ := http.NewRequest(http.MethodPost, "/confirm", bytes.NewBuffer(bodyJSON))
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
