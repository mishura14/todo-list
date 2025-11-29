package servise

import (
	"net"
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

	// Проверка домена
	domain := parts[1]
	if _, err := net.LookupMX(domain); err != nil {
		return false
	}

	return true
}
