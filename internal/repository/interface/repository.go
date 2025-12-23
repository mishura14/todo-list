package repository

import (
	"context"
	"git-register-project/internal/models"
	"time"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks
type UserRegister interface {
	CheckEmailExists(email string) (bool, error)
	CreateUser(user *models.UserRedis) error
}

type EmailSender interface {
	SendConfremRegister(toEmail, code string) error
}

type RedisClient interface {
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Del(ctx context.Context, key string) error
}
type UserLogin interface {
	SelectUser(ctx context.Context, email string) (*models.User, error)
	InsertRefreshToken(userID int, refreshTokenHash string) error
}

type TokenGenerator interface {
	CheckHash(password, hash string) bool
	RefreshJWT(userID int) (string, error)
	AccessJWT(userID int) (string, error)
	HashRefreshToken(token string) string
	ValidateToken(tokenStr string) (map[string]interface{}, error)
}
type Refreshtoken interface {
	SelectRefreshToken(ctx context.Context, id int) (string, error)
	UpdateRefreshToken(ctx context.Context, id int, token string, expiresAt time.Time) error
}
