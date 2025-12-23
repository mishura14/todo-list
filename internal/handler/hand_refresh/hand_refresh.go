package handrefresh

import (
	"context"
	"git-register-project/internal/models"
	refreshtoken "git-register-project/internal/servise/refreshToken"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerRefresh struct {
	service *refreshtoken.RefreshService
}

func NewHandlerRefresh(service *refreshtoken.RefreshService) *HandlerRefresh {
	return &HandlerRefresh{
		service: service,
	}
}

func (h *HandlerRefresh) Refresh(c *gin.Context) {
	ctx := context.Background()
	var token models.Refresh
	if err := c.BindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}
	access, refresh, err := h.service.UpdateRefreshToken(ctx, token)
	if err != nil {
		switch err {
		case refreshtoken.ErrInvalidateToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка при валидации токена"})
			return
		case refreshtoken.ErrSelectToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при запросе токена в бд"})
			return
		case refreshtoken.ErrHashToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "недействительный токен"})
			return
		case refreshtoken.ErrGenerateAccess:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при генерации access токена"})
			return
		case refreshtoken.ErrGenerateRefresh:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при генерации refresh токена"})
			return
		case refreshtoken.ErrUpdateToken:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обновлении токена"})
			return

		}
	}
	c.JSON(http.StatusOK, gin.H{
		"refresh_token": refresh,
		"access_token":  access,
	})
}
