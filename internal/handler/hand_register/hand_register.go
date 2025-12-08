package handler_register

import (
	"git-register-project/internal/models"
	"git-register-project/internal/servise/register"
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
			c.JSON(http.StatusBadGateway, gin.H{"error": "email уже зарегистрирован"})
			return
		case register.ErrBadPasswordFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат пароля"})
			return
		case register.ErrHashPassword:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка хеширования пароля"})
			return
		case register.ErrSendConfirmation:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка отправки кода подтверждения"})
			return
		case register.ErrSaveRedis:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения в Redis"})
			return
		case register.ErrSerializeUser:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации пользователя"})
			return
		case register.ErrCheckEmailInDB:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка проверки email в базе данных"})
			return

		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "код регистрации отправлен на email (активен 5 минут)"})
}
