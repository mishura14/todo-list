package handler_register

import (
	"encoding/json"
	"git-register-project/internal/Database/redis"
	"git-register-project/internal/models"
	"git-register-project/internal/repository"
	serversmtp "git-register-project/internal/server_smtp"
	"git-register-project/internal/servise"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные формат данных"})
		return
	}
	exec := servise.CheckEmail(user.Email)
	if !exec {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат email"})
		return
	}
	var exists bool
	exists, err := repository.CheckEmailExists(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неудалось проверить email в базе данных"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email уже зарегистрирован"})
		return
	}
	if servise.CheckPassword(user.Password) {
		hashpassword, err := servise.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка хеширования пароля"})
			return
		}
		code := servise.GenerateSecureCode()
		if err := serversmtp.SendConfremRegister(user.Email, code); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		user_redis := models.UserRedis{
			Name:     user.Name,
			Email:    user.Email,
			Password: hashpassword,
			Code:     code,
		}
		user_redisJSON, err := json.Marshal(user_redis)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации данных"})
			return
		}
		//сохраняем данные в redis
		err = redis.RDB.Set(redis.Ctx, code, user_redisJSON, 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения данных в redis"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "код регистрации отправлен на email, код станет не действителен через 5 минут"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пароль не соответствует требованиям"})
	}

}
