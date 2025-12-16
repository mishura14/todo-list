package redis

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis хранит клиент Redis
type Redis struct {
	Client *redis.Client
	Ctx    context.Context
}

// ConnectRedis подключается к Redis, без запуска контейнера
func ConnectRedis() (*Redis, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // значение по умолчанию
	}

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // пароль по умолчанию
		DB:       0,  // default DB
	})

	// Проверяем подключение
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Client: client,
		Ctx:    ctx,
	}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return r.Client.Get(ctx, key).Bytes()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
