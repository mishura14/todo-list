package handler_register

import (
	"git-register-project/internal/models"
	"git-register-project/internal/servise/register/register"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerRegister struct {
	service *register.RegisterService
}

func NewRegister(service *register.RegisterService) *HandlerRegister {
	return &HandlerRegister{
		service: service,
	}
}

func (h *HandlerRegister) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}
	err := h.service.Register(&user)
	if err != nil {
		switch err {
		case register.ErrBadEmailFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат email"})
			return
		case register.ErrEmailExists:
			c.JSON(http.StatusConflict, gin.H{"error": "email уже зарегистрирован"})
			return
		case register.ErrBadPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": "пароль не соответствует требованиям"})
			return
		case register.ErrHashPassword:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка хеширования пароля"})
			return
		case register.ErrSendConfirmation:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось отправить письмо"})
			return
		case register.ErrSaveRedis:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения в redis"})
			return
		case register.ErrSerializeUser:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации пользовател"})
			return
		case register.ErrCheckEmailInDB:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки email в базе"})
			return

		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "код регистрации отправлен на email (активен 5 минут)"})
}
