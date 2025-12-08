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
