package main

import (
	"fmt"
	"git-register-project/internal/Database/postgres"
	"git-register-project/internal/Database/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//подключение  postgresql
	db, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Database connection established")
	//подключение redis
	rdb, err := redis.ConnectRedis()
	if err != nil {
		panic(err)
	}
	defer rdb.Close()
	fmt.Println("Connection Redis")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
