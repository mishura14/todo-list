package refreshtoken

import (
	"context"
	"git-register-project/internal/models"
	"git-register-project/internal/repository/interface/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		setupMocks    func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken)
		token         models.Refresh
		expectAccess  string
		expectRefresh string
		expected      error
	}{
		{
			name:  "success",
			token: models.Refresh{Token: "valid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(1)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), 1).Return("hash-refresh", nil)
				tg.EXPECT().CheckHash("valid refresh", "hash-refresh").Return(false)
				tg.EXPECT().AccessJWT(1).Return("new-Access", nil)
				tg.EXPECT().RefreshJWT(1).Return("new-Refresh", nil)
				tg.EXPECT().HashRefreshToken("new-Refresh").Return("new-hash")
				repo.EXPECT().UpdateRefreshToken(gomock.Any(), 1, "new-hash", gomock.AssignableToTypeOf(time.Time{})).Return(nil)
			},
			expectAccess:  "new-Access",
			expectRefresh: "new-Refresh",
			expected:      nil,
		},
		{
			name:  "invalid token",
			token: models.Refresh{Token: "invalid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("invalid refresh").Return(nil, ErrInvalidateToken)
			},
			expected: ErrInvalidateToken,
		},
		{
			name:  "err select token",
			token: models.Refresh{Token: "valid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(1)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), 1).Return("", ErrSelectToken)
			},
			expected: ErrSelectToken,
		},
		{
			name:  "err hash token",
			token: models.Refresh{Token: "valid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(1)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), 1).Return("hash-refresh", nil)
				tg.EXPECT().CheckHash("valid refresh", "hash-refresh").Return(false)
			},
			expected: ErrHashToken,
		},
		{
			name:  "err generate Access",
			token: models.Refresh{Token: "valid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(1)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), 1).Return("hash-refresh", nil)
				tg.EXPECT().CheckHash("valid refresh", "hash-refresh").Return(true)
				tg.EXPECT().AccessJWT(1).Return("", ErrGenerateAccess)
			},
			expected: ErrGenerateAccess,
		},
		{
			name:  "err generate Refresh",
			token: models.Refresh{Token: "valid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(1)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), 1).Return("hash-refresh", nil)
				tg.EXPECT().CheckHash("valid refresh", "hash-refresh").Return(true)
				tg.EXPECT().AccessJWT(1).Return("new-Access", nil)
				tg.EXPECT().RefreshJWT(1).Return("", ErrGenerateRefresh)
			},
			expected: ErrGenerateRefresh,
		},
		{
			name:  "err update token",
			token: models.Refresh{Token: "valid refresh"},
			setupMocks: func(tg *mocks.MockTokenGenerator, repo *mocks.MockRefreshtoken) {
				tg.EXPECT().ValidateToken("valid refresh").Return(map[string]interface{}{"user_id": float64(1)}, nil)
				repo.EXPECT().SelectRefreshToken(gomock.Any(), 1).Return("hash-refresh", nil)
				tg.EXPECT().CheckHash("valid refresh", "hash-refresh").Return(true)
				tg.EXPECT().AccessJWT(1).Return("new-Access", nil)
				tg.EXPECT().RefreshJWT(1).Return("new-Refresh", nil)
				tg.EXPECT().HashRefreshToken("new-Refresh").Return("new-hash")
				repo.EXPECT().UpdateRefreshToken(gomock.Any(), 1, "new-hash", gomock.AssignableToTypeOf(time.Time{})).Return(ErrUpdateToken)
			},
			expected: ErrUpdateToken,
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
				tt.setupMocks(tg, repo)
			}
			service := NewRefreshService(tg, redis, repo)
			access, refresh, err := service.UpdateRefreshToken(ctx, tt.token)
			if tt.expected != nil {
				require.Error(t, err)
				require.Equal(t, tt.expected, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectAccess, access)
			require.Equal(t, tt.expectRefresh, refresh)
		})
	}
}
