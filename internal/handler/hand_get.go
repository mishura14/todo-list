package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}
