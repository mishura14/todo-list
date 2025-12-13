package generatejwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// функция для генерации JWT токена
func GenerateJWT(userID int) (string, error) {
	//получаем секретный ключ из переменной окружения
	secret := os.Getenv("JWT_SECRET")
	//проверяем, что секретный ключ не пустой
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable not set")
	}
	//создаем JWT токен
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	//
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
