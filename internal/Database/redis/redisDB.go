package redis

import (
	"context"
	"os"
	"os/exec"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

// ConnectRedis подключение к redis и запуск если не запущен
func ConnectRedis() (*redis.Client, error) {
	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		redis_port = "localhost:6379" // default port
	}

	//Проверяем подключение
	client := redis.NewClient(&redis.Options{
		Addr:     redis_port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//проверяем подключение
	_, err := client.Ping(Ctx).Result()
	if err == nil {
		RDB = client // Инициализируем глобальную переменную
		return client, nil
	}

	//если не удалось подключиться, запускаем redis docker
	if err := startRedis(); err != nil {
		return nil, err
	}

	// подключаемся снова после запуска Redis
	_, err = client.Ping(Ctx).Result()
	if err != nil {
		return nil, err
	}

	RDB = client // Инициализируем глобальную переменную
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
