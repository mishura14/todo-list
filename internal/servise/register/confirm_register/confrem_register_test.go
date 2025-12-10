package comfirm_register

import (
	"errors"
	"git-register-project/internal/models"
	"git-register-project/internal/repository/mocks"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestConfirmRegister(t *testing.T) {
	cases := []struct {
		name        string
		code        string
		setupMocks  func(repo *mocks.MockUserRegister, redis *mocks.MockRedisClient)
		expectedErr error
	}{
		{
			name: "success",
			code: "123456",
			setupMocks: func(repo *mocks.MockUserRegister, redis *mocks.MockRedisClient) {
				data := []byte(`{"Name":"John","Email":"john@example.com","Password":"hash","Code":"123456"}`)
				redis.EXPECT().Get(gomock.Any(), "123456").Return(data, nil)
				repo.EXPECT().CreateUser(&models.UserRedis{Name: "John", Email: "john@example.com", Password: "hash"}).Return(nil)
				redis.EXPECT().Del(gomock.Any(), "123456").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "Redis key missing",
			code: "123456",
			setupMocks: func(repo *mocks.MockUserRegister, redis *mocks.MockRedisClient) {
				redis.EXPECT().Get(gomock.Any(), "123456").Return(nil, errors.New("not found"))
			},
			expectedErr: ErrCodeTimeout,
		},
		{
			name: "Bad JSON",
			code: "123456",
			setupMocks: func(repo *mocks.MockUserRegister, redis *mocks.MockRedisClient) {
				redis.EXPECT().Get(gomock.Any(), "123456").Return([]byte("bad json"), nil)
			},
			expectedErr: ErrBadJSONFormat,
		},
		{
			name: "CreateUser error",
			code: "123456",
			setupMocks: func(repo *mocks.MockUserRegister, redis *mocks.MockRedisClient) {
				data := []byte(`{"Name":"John","Email":"john@example.com","Password":"hash","Code":"123456"}`)
				redis.EXPECT().Get(gomock.Any(), "123456").Return(data, nil)
				repo.EXPECT().CreateUser(&models.UserRedis{Name: "John", Email: "john@example.com", Password: "hash"}).Return(errors.New("db error"))
			},
			expectedErr: ErrCreateUser,
		},
		{
			name: "Redis delete error",
			code: "123456",
			setupMocks: func(repo *mocks.MockUserRegister, redis *mocks.MockRedisClient) {
				data := []byte(`{"Name":"John","Email":"john@example.com","Password":"hash","Code":"123456"}`)
				redis.EXPECT().Get(gomock.Any(), "123456").Return(data, nil)
				repo.EXPECT().CreateUser(&models.UserRedis{Name: "John", Email: "john@example.com", Password: "hash"}).Return(nil)
				redis.EXPECT().Del(gomock.Any(), "123456").Return(errors.New("redis delete error"))
			},
			expectedErr: ErrDelCodeUser,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUserRegister(ctrl)
			redis := mocks.NewMockRedisClient(ctrl)

			tt.setupMocks(repo, redis)

			svc := NewConfirmRegisterService(repo, nil, redis)
			err := svc.ConfirmRegister(tt.code)

			if err != tt.expectedErr {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
