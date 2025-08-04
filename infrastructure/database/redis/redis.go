package redis

import (
	"github.com/redis/go-redis/v9"
)

// NewRedis initializes a new Redis client using the provided configuration options.
// It takes a *redis.Options struct as input, which contains the Redis server address,
// password, and database index, and returns a *redis.Client instance.
func NewRedis(config *redis.Options) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return rdb
}
