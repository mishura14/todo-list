package login

import (
	"context"
	"errors"
	"git-register-project/internal/models"
	repository "git-register-project/internal/repository/interface"
)

type LoginService struct {
	login repository.TokenGenerator
	repo  repository.UserLogin
	mail  repository.EmailSender
	redis repository.RedisClient
}

func NewLoginService(repo repository.UserLogin, mail repository.EmailSender, redis repository.RedisClient, login repository.TokenGenerator) *LoginService {
	return &LoginService{
		login: login,
		repo:  repo,
		mail:  mail,
		redis: redis,
	}
}

var (
	ErrUserNotFound       = errors.New("ваш аккаунт не зарегистрирован")
	ErrPasswordIncorrect  = errors.New("неверный пароль")
	ErrAccessToken        = errors.New("ошибка создания токена")
	ErrRefreshToken       = errors.New("ошибка создания refresh токена")
	ErrHashRefreshToken   = errors.New("ошибка хеширования refresh токена")
	ErrInsertRefreshToken = errors.New("ошибка вставки refresh токена")
)

func (ls *LoginService) Login(login models.UserLogin) (accessToken string, refreshToken string, err error) {
	ctx := context.Background()
	//получаем данные по email
	user, err := ls.repo.SelectUser(ctx, login.Email)
	if err != nil {
		return "", "", ErrUserNotFound
	}
	//сверяем пароль
	if !ls.login.CheckPasswordHash(login.Password, user.Password) {
		return "", "", ErrPasswordIncorrect
	}
	//создание токенов
	token, err := ls.login.AccessJWT(user.ID)
	if err != nil {
		return "", "", ErrAccessToken
	}
	refreshToken, err = ls.login.RefreshJWT(user.ID)
	if err != nil {
		return "", "", ErrRefreshToken
	}
	//хеширование токена
	tokenHash, err := ls.login.HashPassword(refreshToken)
	if err != nil {
		return "", "", ErrHashRefreshToken
	}
	//вставка токена в БД
	err = ls.repo.InsertRefreshToken(ctx, user.ID, tokenHash)
	if err != nil {
		return "", "", ErrInsertRefreshToken
	}
	return token, tokenHash, nil

}
