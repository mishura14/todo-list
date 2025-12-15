package useCase

import (
	"context"
	"git-register-project/internal/models"
)

func (r *PostgreUser) SelectUser(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, email, password_hash FROM users WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
