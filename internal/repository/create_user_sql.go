package repository

import (
	"context"
	"git-register-project/internal/Database/postgres"
	"git-register-project/internal/models"
)

// запрос создания пользователя в бд
func CreateUser(user *models.UserRedis) error {
	_, err := postgres.DB.ExecContext(context.Background(),
		`INSERT INTO users (name, email, password_hash) VALUES ($1,$2,$3)`,
		user.Name, user.Email, user.Password)
	return err
}
