package checkhash

import (
	"crypto/subtle"
	hashrefreshtoken "git-register-project/internal/servise/jwt_token/hashRefreshToken"
)

// CompareRefreshToken сравнивает raw-токен и сохранённый hash
func CheckRefreshToken(rawToken string, storedHash string) bool {
	rawHash := hashrefreshtoken.HashRefreshToken(rawToken)

	// constant time compare (защита от timing attack)
	return subtle.ConstantTimeCompare(
		[]byte(rawHash),
		[]byte(storedHash),
	) == 1
}
