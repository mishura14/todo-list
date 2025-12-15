package useCase

import (
	"context"
)

// запрос в базу данных проверки существования email
func (r *PostgreUser) CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}
