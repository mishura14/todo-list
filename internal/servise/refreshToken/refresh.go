package refreshtoken

import (
	"context"
	"errors"
	"git-register-project/internal/models"
	repository "git-register-project/internal/repository/interface"
	"time"
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

var (
	ErrInvalidateToken = errors.New("ошибка при валидации токена")
	ErrSelectToken     = errors.New("ошибка при запросе токена в бд")
	ErrHashToken       = errors.New("недействительный токен")
	ErrGenerateAccess  = errors.New("ошибка при генерации access токена")
	ErrGenerateRefresh = errors.New("ошибка при генерации refresh токена")
	ErrUpdateToken     = errors.New("ошибка при обновлении токена")
)

func (s *RefreshService) UpdateRefreshToken(ctx context.Context, token models.Refresh) (access, refresh string, err error) {
	claims, err := s.repoReft.ValidateToken(token.Token)
	if err != nil {
		return "", "", ErrInvalidateToken
	}
	id := int(claims["user_id"].(float64))
	hash_token, err := s.repo.SelectRefreshToken(ctx, id)
	if err != nil {
		return "", "", ErrSelectToken
	}
	if s.repoReft.CheckHash(token.Token, hash_token) {
		return "", "", ErrHashToken
	}
	newAccess, err := s.repoReft.AccessJWT(id)
	if err != nil {
		return "", "", ErrGenerateAccess
	}
	newRefresh, err := s.repoReft.RefreshJWT(id)
	if err != nil {
		return "", "", ErrGenerateRefresh
	}
	newHashRefresh := s.repoReft.HashRefreshToken(newRefresh)
	err = s.repo.UpdateRefreshToken(ctx, id, newHashRefresh, time.Now().Add(time.Hour*7*24))
	if err != nil {
		return "", "", ErrUpdateToken
	}

	return newAccess, newRefresh, nil
}
