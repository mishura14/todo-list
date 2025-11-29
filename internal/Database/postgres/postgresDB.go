package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() (*sql.DB, error) {

	// Параметры подключения с значениями по умолчанию
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "postgres")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Пытаемся подключиться
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %v", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		fmt.Println("🚀 Запускаем PostgreSQL в Docker...")

		// Запускаем контейнер
		if err := startDB(); err != nil {
			return nil, fmt.Errorf("ошибка запуска БД: %v", err)
		}

		// Ждем запуска
		fmt.Println("⏳ Ожидаем запуска PostgreSQL...")
		time.Sleep(5 * time.Second)

		// Пытаемся снова
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка подключения после запуска: %v", err)
		}

		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("не удалось подключиться к БД: %v", err)
		}
	}

	// Инициализируем глобальную переменную
	DB = db
	return DB, nil
}

func startDB() error {
	// Проверяем есть ли контейнер
	cmd := exec.Command("docker", "inspect", "postgres")
	if cmd.Run() == nil {
		fmt.Println("🔄 Запускаем существующий контейнер PostgreSQL...")
		return exec.Command("docker", "start", "postgres").Run()
	}

	// Создаем новый контейнер
	fmt.Println("📦 Создаем новый контейнер PostgreSQL...")
	cmd = exec.Command("docker", "run", "-d",
		"--name", "postgres",
		"-e", "POSTGRES_PASSWORD=postgres",
		"-e", "POSTGRES_DB=postgres",
		"-p", "5432:5432",
		"postgres:15")
	return cmd.Run()
}

// Вспомогательная функция для получения переменной окружения с значением по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
