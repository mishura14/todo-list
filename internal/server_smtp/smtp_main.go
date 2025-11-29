package serversmtp

import (
	"os"
	"strconv"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

func getSMTPClient() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()

	// Получаем настройки SMTP из переменных окружения
	server.Host = getEnvOrDefault("SMTP_HOST", "smtp.gmail.com")

	port := getEnvOrDefault("SMTP_PORT", "587")
	if p, err := strconv.Atoi(port); err == nil {
		server.Port = p
	} else {
		server.Port = 587
	}

	server.Username = getEnvOrDefault("SMTP_EMAIL", "")
	server.Password = getEnvOrDefault("SMTP_PASSWORD", "")
	server.Encryption = mail.EncryptionSTARTTLS
	server.Authentication = mail.AuthLogin
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	return server.Connect()
}

// Вспомогательная функция для получения переменной окружения с значением по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
