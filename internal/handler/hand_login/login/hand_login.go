package handlogin

import (
	"git-register-project/internal/models"
	"git-register-project/internal/servise/login"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerLogin struct {
	service *login.LoginService
}

func NewLogin(service *login.LoginService) *HandlerLogin {
	return &HandlerLogin{
		service: service,
	}
}
func (h *HandlerLogin) Login(c *gin.Context) {
	var user models.UserLogin
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}
	token, tokenHash, err := h.service.Login(user)
	if err != nil {
		switch err {
		case login.ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ваш аккаунт не зарегистрирован"})
			return
		case login.ErrPasswordIncorrect:
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный пароль"})
			return
		case login.ErrAccessToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания токена"})
			return
		case login.ErrRefreshToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания refresh токена"})
			return
		case login.ErrHashRefreshToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка хеширования refresh токена"})
			return
		case login.ErrInsertRefreshToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка вставки refresh токена"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"token_hash": tokenHash,
	})
}
