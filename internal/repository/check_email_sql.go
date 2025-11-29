package repository

import (
	"context"
	"git-register-project/internal/Database/postgres"
)

// запрос в базу данных проверки существования email
func CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := postgres.DB.QueryRowContext(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}
