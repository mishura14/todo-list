package useCase

import (
	"context"
	"fmt"
	"time"
)

// фунция внесения refresh token в bd
func (r *PostgreUser) InsertRefreshToken(userID int, refreshTokenHash string) error {

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	_, err := r.db.ExecContext(context.Background(),
		`INSERT INTO refresh_tokens(user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (user_id)
		 DO UPDATE SET
		     token_hash = EXCLUDED.token_hash,
		     expires_at = EXCLUDED.expires_at`,
		userID,
		refreshTokenHash,
		expiresAt,
	)
	fmt.Println(err)
	return err
}
