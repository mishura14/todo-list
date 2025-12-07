package handler_register

import (
	"encoding/json"
	"git-register-project/internal/Database/redis"
	"git-register-project/internal/models"
	"git-register-project/internal/repository"
	"git-register-project/internal/servise/generate_code"
	"git-register-project/internal/servise/hash_password/password_hash"
	"git-register-project/internal/servise/validate/valid_email"
	"git-register-project/internal/servise/validate/valid_password"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Register struct {
	redis *redis.Redis
	repo  repository.UserRegister
	mail  repository.EmailSender
}

func NewRegister(redis *redis.Redis, repo repository.UserRegister, mail repository.EmailSender) *Register {
	return &Register{
		redis: redis,
		repo:  repo,
		mail:  mail,
	}
}

// оброботчик запросов регистрации
func (r *Register) Register(c *gin.Context) {
	var user models.User
	//проверка получения данных
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверные формат данных"})
		return
	}
	//проверка формата email
	exec := valid_email.CheckEmail(user.Email)
	if !exec {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат email"})
		return
	}
	var exists bool
	//проверка email на существование в базе данных
	exists, err := r.repo.CheckEmailExists(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неудалось проверить email в базе данных"})
		return
	}
	//проверка email на существование в базе данных
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email уже зарегистрирован"})
		return
	}
	//проверка пароля на соответствие требованиям и хеширование
	if valid_password.CheckPassword(user.Password) {
		hashpassword, err := password_hash.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка хеширования пароля"})
			return
		}
		//генерация кода подтверждения
		code := generate_code.GenerateSecureCode()
		if err := r.mail.SendConfremRegister(user.Email, code); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		user_redis := models.UserRedis{
			Name:     user.Name,
			Email:    user.Email,
			Password: hashpassword,
			Code:     code,
		}
		//преобразование в json
		user_redisJSON, err := json.Marshal(user_redis)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации данных"})
			return
		}
		//сохраняем данные в redis
		err = r.redis.Client.Set(r.redis.Ctx, code, user_redisJSON, 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения данных в redis"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "код регистрации отправлен на email, код станет не действителен через 5 минут"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пароль не соответствует требованиям"})
	}

}
