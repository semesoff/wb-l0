package cache

import (
	"github.com/go-redis/redis"
	"wb-l0/config"
)

type Redis struct {
	client *redis.Client
}

func InitRedisClient() *Redis {
	cfg := config.GetConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
	})
	return &Redis{client: client}
}

func (r *Redis) AddCache(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) GetCache(key string) (string, error) {
	return r.client.Get(key).Result()
}
