package validatejwt

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateToken проверяет JWT access-токен и возвращает claims
func ValidateToken(tokenStr string) (map[string]interface{}, error) {
	// Получаем секретный ключ из переменной окружения
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET не найден")
	}

	// Разбираем и проверяем токен
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("не поддерживаемый алгоритм подписи")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("неверный токен")
	}

	// Проверяем, что токен валиден
	if !token.Valid {
		return nil, errors.New("неверный токен")
	}

	// Получаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("не удалось получить claims")
	}

	return claims, nil
}
