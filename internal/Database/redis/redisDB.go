package redis

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ConnectResid  подключение к redis и запусе если не запущен
func ConnectRedis() (*redis.Client, error) {
	errd := godotenv.Load("/home/mishura/ZedProject/git-register-project/.env")
	if errd != nil {
		fmt.Println("Error loading .env file:", errd)
	}
	redis_port := os.Getenv("REDIS_PORT")

	//Проверям продключение
	client := redis.NewClient(&redis.Options{
		Addr:     redis_port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	//проверяем подключение
	_, err := client.Ping(ctx).Result()
	if err == nil {
		return client, nil
	}
	//если не удолось запускаем redis docker
	if err := startRedis(); err != nil {
		return nil, err
	}
	// подключаемся снова
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
func startRedis() error {
	cmd := exec.Command("docker", "inspect", "redis")
	if cmd.Run() == nil {
		return exec.Command("docker", "start", "redis").Run()
	}
	cmd = exec.Command(
		"docker", "run", "-d",
		"--name", "redis",
		"-p", "6379:6379",
		"redis:7-alpine",
	)
	return cmd.Run()
}
