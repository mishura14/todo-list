package handler_register

import (
	"encoding/json"
	"git-register-project/internal/models"
	"git-register-project/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// обработчик подтверждения регистрации
func (r *Redis) Confirm_register(c *gin.Context) {
	var code models.CodeUser
	//получение и обработка кода подтверждения регистрации
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат кода"})
		return
	}
	//поиск кода в Redis
	val, err := r.redisClient.Client.Get(r.redisClient.Ctx, code.Code).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "код не найден или время истекло"})
		return
	}
	//приведение полученных данных к типу UserRedis
	var user_from_redis models.UserRedis
	if err := json.Unmarshal([]byte(val), &user_from_redis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при десериализации данных"})
		return
	}
	//регистрация пользователя
	err = repository.CreateUser(&user_from_redis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании пользователя"})
		return
	}
	//удаление кода из Redis
	err = r.redisClient.Client.Del(r.redisClient.Ctx, code.Code).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении кода"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "регистрация успешно завершена"})
}
