package useCase

import "context"

func (r *PostgreUser) SelectRefreshToken(ctx context.Context, id int) (string, error) {
	var refreshToken string
	err := r.db.QueryRowContext(
		ctx,
		"SELECT refresh_token FROM refresh_tokens WHERE user_id = $1",
		id,
	).Scan(&refreshToken)
	if err != nil {
		return "", err
	}
	return refreshToken, err
}
