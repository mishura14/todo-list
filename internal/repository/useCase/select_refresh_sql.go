package useCase

import "context"

func (r *PostgreRefresh) SelectRefreshToken(ctx context.Context, id int) (string, error) {
	var refreshToken string
	err := r.db.QueryRowContext(
		ctx,
		"SELECT token_hash FROM refresh_tokens WHERE user_id = $1",
		id,
	).Scan(&refreshToken)
	if err != nil {
		return "", err
	}
	return refreshToken, err
}
