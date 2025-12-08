package register

import (
	"context"
	"encoding/json"
	"errors"
	"git-register-project/internal/models"
)

var (
	ErrBadJSONFormat = errors.New("ошибка формата JSON")
	ErrCodeTimeout   = errors.New("время кода истекло")
	ErrCreateUser    = errors.New("ошибка создания пользователя")
	ErrDelCodeUser   = errors.New("ошибка удаления кода пользователя")
)

func (s *RegisterService) ConfirmRegister(code string) error {
	ctx := context.Background()
	val, err := s.redis.Get(ctx, code)
	if err != nil {
		return ErrCodeTimeout
	}
	var user_from_redis models.UserRedis
	if err := json.Unmarshal([]byte(val), &user_from_redis); err != nil {
		return ErrBadJSONFormat
	}
	err = s.repo.CreateUser(&models.UserRedis{
		Name:     user_from_redis.Name,
		Email:    user_from_redis.Email,
		Password: user_from_redis.Password,
	})
	if err != nil {
		return ErrCreateUser
	}
	err = s.redis.Del(ctx, code)
	if err != nil {
		return ErrDelCodeUser
	}
	return nil
}
