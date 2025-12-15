package useCase

import (
	"context"
	"time"
)

func (r *PostgreUser) InsertRefreshToken(ctx context.Context, userID int, refreshTokenHash string) error {

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id)
		 DO UPDATE SET
		     token_hash = EXCLUDED.token_hash,
		     expires_at = EXCLUDED.expires_at`,
		userID,
		refreshTokenHash,
		expiresAt,
	)

	return err
}
