package jwt_adapter

import (
	checkhash "git-register-project/internal/servise/bcryptHash/checkHash"
	accessjwt "git-register-project/internal/servise/jwt_token/accessJWT"
	hashrefreshtoken "git-register-project/internal/servise/jwt_token/hashRefreshToken"
	refreshjwt "git-register-project/internal/servise/jwt_token/refresh_token"
)

// Adapter реализует интерфейс TokenGenerator
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) AccessJWT(userID int) (string, error) {
	return accessjwt.AccessJWT(userID)
}

func (a *Adapter) RefreshJWT(userID int) (string, error) {
	return refreshjwt.RefreshJWT(userID)
}

func (a *Adapter) HashRefreshToken(token string) string {
	return hashrefreshtoken.HashRefreshToken(token)
}

func (a *Adapter) CheckHash(password, hash string) bool {
	return checkhash.CheckHash(password, hash)
}
