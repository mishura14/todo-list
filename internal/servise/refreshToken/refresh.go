package refreshtoken

import (
	"context"
	repository "git-register-project/internal/repository/interface"
)

type RefreshService struct {
	repoReft repository.TokenGenerator
	redis    repository.RedisClient
	repo     repository.Refreshtoken
}

func NewRefreshService(repoReft repository.TokenGenerator, redis repository.RedisClient, repo repository.Refreshtoken) *RefreshService {
	return &RefreshService{
		repoReft: repoReft,
		redis:    redis,
		repo:     repo,
	}
}
func (s *RefreshService) UpdateRefreshToken(ctx context.Context, id int, token string) (access, refresh string, err error) {

	return access, refresh, err
}
