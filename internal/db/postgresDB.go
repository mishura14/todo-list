package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Проверяет работает ли БД
func isDBRunning() bool {
	errd := godotenv.Load("/home/mishura/ZedProject/git-register-project/.env")
	if errd != nil {
		fmt.Println("Error loading .env file:", errd)
	}
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")

	db, err := sql.Open("", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, db_host, db_port, db_name))
	if err != nil {
		return false
	}
	defer db.Close()

	return db.Ping() == nil
}

// Подключается к БД
func connectToDB() (*sql.DB, error) {
	errd := godotenv.Load("/home/mishura/ZedProject/git-register-project/.env")
	if errd != nil {
		fmt.Println("Error loading .env file:", errd)
	}
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db_user, db_password, db_host, db_port, db_name))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Запускает контейнер с БД
func startDB() error {
	// Проверяем есть ли контейнер
	cmd := exec.Command("docker", "inspect", "postgres")
	if cmd.Run() == nil {
		// Контейнер есть, запускаем
		return exec.Command("docker", "start", "postgres").Run()
	}

	// Создаем новый контейнер
	cmd = exec.Command("docker", "run", "-d",
		"--name", "postgres",
		"-e", "POSTGRES_PASSWORD=postgres",
		"-p", "5432:5432",
		"postgres:15")

	return cmd.Run()
}

// функция подключения к БД
func ConnectDB() (*sql.DB, error) {
	// Если БД уже работает - подключаемся
	if isDBRunning() {
		return connectToDB()
	}

	// Запускаем БД
	if err := startDB(); err != nil {
		return nil, fmt.Errorf("ошибка запуска БД: %v", err)
	}

	// Ждем запуска БД
	time.Sleep(5 * time.Second)

	// Пытаемся подключиться
	return connectToDB()
}
