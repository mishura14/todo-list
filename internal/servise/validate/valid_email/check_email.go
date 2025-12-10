package valid_email

import (
	"regexp"
	"strings"
)

func CheckEmail(email string) bool {
	// Проверка длины
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Проверка формата
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		return false
	}

	// Проверка что есть домен
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return true
}
