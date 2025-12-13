package validatejwt

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// функция для проверки токена
func ValidateToken(tokenstr string) (map[string]interface{}, error) {
	//получаем секретный ключ из переменной окружения
	secret := os.Getenv("JWT_SECRET")
	//проверяем, что секретный ключ не пустой
	if secret == "" {
		return nil, errors.New("JWT_SECRET не найдет")
	}
	//создаем функцию для проверки токена
	tokens, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("не поддерживаемый алгоритм")
		}
		return []byte(secret), nil
	})
	//проверяем, что токен не пустой и валиден
	if err != nil || !tokens.Valid {
		return nil, errors.New("неверный токен")
	}
	claims, ok := tokens.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("не получилось получить токен")
	}
	return claims, nil
}
