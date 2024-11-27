package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"wb-l0/config"
	"wb-l0/internal/models/order"
)

type RedisProvider interface {
	AddCache(key string, value interface{}) error
	GetCache(key string) ([]byte, error)
	BytesToModel(key string) (order.Order, error)
}

type Redis struct {
	client *redis.Client
}

func NewRedis(cfg *config.Redis) *Redis {
	r := &Redis{}
	r.InitRedisClient(cfg)
	return r
}

func (r *Redis) InitRedisClient(cfg *config.Redis) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password,
	})
	r.client = client
}

func (r *Redis) AddCache(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) GetCache(key string) ([]byte, error) {
	return r.client.Get(key).Bytes()
}

func (r *Redis) BytesToModel(key string) (order.Order, error) {
	data, err := r.GetCache(key)
	if err != nil {
		return order.Order{}, err
	}

	var orderData order.Order
	if err := json.Unmarshal(data, &orderData); err != nil {
		return order.Order{}, err
	}
	fmt.Println(orderData)
	return orderData, nil
}

func (r *Redis) RestoreCacheFromDB() {

}
