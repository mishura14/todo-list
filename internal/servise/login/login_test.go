package login

import (
	"context"
	"errors"
	"git-register-project/internal/models"
	"git-register-project/internal/repository/interface/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLoginService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserLogin(ctrl)
	mockMail := mocks.NewMockEmailSender(ctrl)
	mockRedis := mocks.NewMockRedisClient(ctrl)
	mockToken := mocks.NewMockTokenGenerator(ctrl)

	service := NewLoginService(mockRepo, mockMail, mockRedis, mockToken)

	tests := []struct {
		name          string
		input         models.UserLogin
		setupMocks    func()
		expectedError error
	}{
		{
			name: "success",
			input: models.UserLogin{
				Email:    "ok@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				ctx := context.Background()
				mockRepo.EXPECT().SelectUser(ctx, "ok@example.com").Return(&models.User{
					ID:       1,
					Email:    "ok@example.com",
					Password: "hashedpassword",
				}, nil)
				mockToken.EXPECT().CheckPasswordHash("password123", "hashedpassword").Return(true)
				mockToken.EXPECT().AccessJWT(1).Return("access-token", nil)
				mockToken.EXPECT().RefreshJWT(1).Return("refresh-token", nil)
				mockToken.EXPECT().HashPassword("refresh-token").Return("hashed-refresh", nil)
				mockRepo.EXPECT().InsertRefreshToken(ctx, 1, "hashed-refresh").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			input: models.UserLogin{
				Email:    "no@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				ctx := context.Background()
				mockRepo.EXPECT().SelectUser(ctx, "no@example.com").Return(nil, errors.New("not found"))
			},
			expectedError: ErrUserNotFound,
		},
		{
			name: "wrong password",
			input: models.UserLogin{
				Email:    "ok@example.com",
				Password: "wrongpass",
			},
			setupMocks: func() {
				ctx := context.Background()
				mockRepo.EXPECT().SelectUser(ctx, "ok@example.com").Return(&models.User{
					ID:       1,
					Email:    "ok@example.com",
					Password: "hashedpassword",
				}, nil)
				mockToken.EXPECT().CheckPasswordHash("wrongpass", "hashedpassword").Return(false)
			},
			expectedError: ErrPasswordIncorrect,
		},
		{
			name: "insert refresh token fails",
			input: models.UserLogin{
				Email:    "ok@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				ctx := context.Background()
				mockRepo.EXPECT().SelectUser(ctx, "ok@example.com").Return(&models.User{
					ID:       1,
					Email:    "ok@example.com",
					Password: "hashedpassword",
				}, nil)
				mockToken.EXPECT().CheckPasswordHash("password123", "hashedpassword").Return(true)
				mockToken.EXPECT().AccessJWT(1).Return("access-token", nil)
				mockToken.EXPECT().RefreshJWT(1).Return("refresh-token", nil)
				mockToken.EXPECT().HashPassword("refresh-token").Return("hashed-refresh", nil)
				mockRepo.EXPECT().InsertRefreshToken(ctx, 1, "hashed-refresh").Return(errors.New("db error"))
			},
			expectedError: ErrInsertRefreshToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			access, refresh, err := service.Login(tt.input)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Empty(t, access)
				assert.Empty(t, refresh)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "access-token", access)
				assert.Equal(t, "hashed-refresh", refresh)
			}
		})
	}
}
