package register

import (
	"errors"
	"testing"

	"git-register-project/internal/models"
	"git-register-project/internal/repository/interface/mocks"

	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		user        *models.User
		setupMocks  func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient)
		expectedErr error
	}{
		{
			name:        "invalid email",
			user:        &models.User{Email: "bad", Password: "StrongPass123!"},
			setupMocks:  func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {},
			expectedErr: ErrBadEmailFormat,
		},
		{
			name: "DB error",
			user: &models.User{Email: "ok@example.com", Password: "StrongPass123!"},
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, errors.New("db err"))
			},
			expectedErr: ErrCheckEmailInDB,
		},
		{
			name: "Mail sending error",
			user: &models.User{Email: "ok@example.com", Password: "StrongPass123!"},
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				mail.EXPECT().SendConfremRegister("ok@example.com", gomock.Any()).Return(errors.New("smtp error"))
			},
			expectedErr: ErrSendConfirmation,
		},
		{
			name: "Redis error",
			user: &models.User{Email: "ok@example.com", Password: "StrongPass123!"},
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				mail.EXPECT().SendConfremRegister("ok@example.com", gomock.Any()).Return(nil)
				redis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("redis error"))
			},
			expectedErr: ErrSaveRedis,
		},
		{
			name: "email exists",
			user: &models.User{Email: "exist@example.com", Password: "StrongPass123!"},
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("exist@example.com").Return(true, nil)
			},
			expectedErr: ErrEmailExists,
		},
		{
			name: "success",
			user: &models.User{Email: "ok@example.com", Password: "StrongPass123!"},
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				mail.EXPECT().SendConfremRegister("ok@example.com", gomock.Any()).Return(nil)
				redis.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name: "invalid password",
			user: &models.User{Email: "ok@example.com", Password: "mishura14"},
			setupMocks: func(repo *mocks.MockUserRegister, mail *mocks.MockEmailSender, redis *mocks.MockRedisClient) {
				repo.EXPECT().CheckEmailExists("ok@example.com").Return(false, nil)
				// mail/redis НЕ должны вызываться!
			},
			expectedErr: ErrBadPasswordFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mocks.NewMockUserRegister(ctrl)
			mail := mocks.NewMockEmailSender(ctrl)
			redis := mocks.NewMockRedisClient(ctrl)

			tt.setupMocks(repo, mail, redis)

			svc := NewRegisterService(repo, mail, redis)
			err := svc.Register(tt.user)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected %v, got %v", tt.expectedErr, err)
			}
		})
	}

}
