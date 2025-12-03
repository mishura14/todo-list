package redis

import (
	"context"
	"os"
	"os/exec"

	"github.com/redis/go-redis/v9"
)

// Redis хранит клиент Redis
type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

// ConnectRedis подключается к redis и запускает если не запущен
func ConnectRedis() (*Redis, error) {
	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		redis_port = "localhost:6379" // default port
	}

	// Создаем контекст
	ctx := context.Background()

	// Проверяем подключение
	client := redis.NewClient(&redis.Options{
		Addr:     redis_port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// проверяем подключение
	_, err := client.Ping(ctx).Result()
	if err == nil {
		return &Redis{
			Client: client,
			Ctx:    ctx,
		}, nil
	}

	// если не удалось подключиться, запускаем redis docker
	if err := startRedis(); err != nil {
		return nil, err
	}

	// подключаемся снова после запуска Redis
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Client: client,
		Ctx:    ctx,
	}, nil
}

// startRedis запуск Redis в Docker
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
