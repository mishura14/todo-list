package main

import (
	"fmt"
	"git-register-project/internal/Database/postgres"
	"git-register-project/internal/Database/redis"
	"git-register-project/internal/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("⚠️  .env файл не найден: %v", err)
	} else {
		fmt.Println("Connected .env")
	}
	//подключение  postgresql
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected postgres")

	//подключение redis
	rdb, err := redis.ConnectRedis()
	if err != nil {
		panic(err)
	}
	defer rdb.Close()
	fmt.Println("Connection Redis")

	r := gin.Default()
	router.SetupRouter(r)

	// Получаем порт из .env или используем по умолчанию
	port := os.Getenv("APP_PORT")
	r.Run(":" + port)
}
