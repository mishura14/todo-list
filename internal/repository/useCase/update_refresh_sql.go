package useCase

import (
	"context"
	"time"
)

func (r *PostgreUser) UpdateRefreshToken(ctx context.Context, id int, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		"UPDATE refresh_tokens SET token_hash = $1, expires_at = $2,created_at = $3 WHERE user_id = $4",
		token,
		expiresAt,
		time.Now(),
		id,
	)
	return err
}
