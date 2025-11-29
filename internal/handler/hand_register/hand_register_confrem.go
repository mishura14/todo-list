package handler_register

import (
	"encoding/json"
	"git-register-project/internal/Database/redis"
	"git-register-project/internal/models"
	"git-register-project/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Confirm_register(c *gin.Context) {
	var code models.CodeUser
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат кода"})
		return
	}
	val, err := redis.RDB.Get(redis.Ctx, code.Code).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "код не найден или время истекло"})
		return
	}
	var user_from_redis models.UserRedis
	if err := json.Unmarshal([]byte(val), &user_from_redis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при десериализации данных"})
		return
	}
	err = repository.CreateUser(&user_from_redis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании пользователя"})
		return
	}
	err = redis.RDB.Del(redis.Ctx, code.Code).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении кода"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "регистрация успешно завершена"})
}
