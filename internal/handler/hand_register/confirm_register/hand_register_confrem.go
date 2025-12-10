package handler_comfirm_register

import (
	"git-register-project/internal/models"
	comfirm_register "git-register-project/internal/servise/register/confirm_register"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerConfirmRegister struct {
	service *comfirm_register.ConfirmRegisterService
}

func NewConfirmRegister(service *comfirm_register.ConfirmRegisterService) *HandlerConfirmRegister {
	return &HandlerConfirmRegister{
		service: service,
	}
}

// обработчик подтверждения регистрации
func (h *HandlerConfirmRegister) Confirm_register(c *gin.Context) {
	var code models.CodeUser
	//получение и обработка кода подтверждения регистрации
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат кода"})
		return
	}
	err := h.service.ConfirmRegister(code.Code)
	if err != nil {
		switch err {
		case comfirm_register.ErrBadJSONFormat:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки JSON"})
			return
		case comfirm_register.ErrCodeTimeout:
			c.JSON(http.StatusBadRequest, gin.H{"error": "код подтверждения устарел"})
			return
		case comfirm_register.ErrCreateUser:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка создания пользователя"})
			return
		case comfirm_register.ErrDelCodeUser:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка удаления кода подтверждения"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "регистрация успешно завершена"})
}
