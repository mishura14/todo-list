package handler_register

import (
	"git-register-project/internal/models"
	"git-register-project/internal/servise/register"
	"net/http"

	"github.com/gin-gonic/gin"
)

// обработчик подтверждения регистрации
func (h *HandlerRegister) Confirm_register(c *gin.Context) {
	var code models.CodeUser
	//получение и обработка кода подтверждения регистрации
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат кода"})
		return
	}
	err := h.service.ConfirmRegister(code.Code)
	if err != nil {
		switch err {
		case register.ErrBadJSONFormat:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки JSON"})
			return
		case register.ErrCodeTimeout:
			c.JSON(http.StatusBadRequest, gin.H{"error": "код подтверждения устарел"})
			return
		case register.ErrCreateUser:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания пользователя"})
			return
		case register.ErrDelCodeUser:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления кода подтверждения"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "регистрация успешно завершена"})
}
