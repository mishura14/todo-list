package handler_register

import (
	"encoding/json"
	"git-register-project/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// обработчик подтверждения регистрации
func (r *Register) Confirm_register(c *gin.Context) {
	var code models.CodeUser
	//получение и обработка кода подтверждения регистрации
	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат кода"})
		return
	}
	//поиск кода в Redis
	val, err := r.redis.Client.Get(r.redis.Ctx, code.Code).Result()
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
	err = r.repo.CreateUser(&models.UserRedis{
		Name:     user_from_redis.Name,
		Email:    user_from_redis.Email,
		Password: user_from_redis.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании пользователя"})
		return
	}
	//удаление кода из Redis
	err = r.redis.Client.Del(r.redis.Ctx, code.Code).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении кода"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "регистрация успешно завершена"})
}
