package repository

import (
	"context"
	"git-register-project/internal/models"
)

// запрос создания пользователя в бд
func (db *Database) CreateUser(user *models.UserRedis) error {
	_, err := db.DB.ExecContext(context.Background(),
		`INSERT INTO users (name, email, password_hash) VALUES ($1,$2,$3)`,
		user.Name, user.Email, user.Password)
	return err
}
