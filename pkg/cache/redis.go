package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisClient "github.com/OmidRasouli/weather-api/infrastructure/database/redis"
	"github.com/redis/go-redis/v9"
)

// RedisCache provides caching operations using Redis
type RedisCache struct {
	redis *redisClient.Redis
	ttl   time.Duration
}

// NewRedisCache creates a new Redis cache service
func NewRedisCache(redis *redisClient.Redis) *RedisCache {
	ttl := time.Duration(redis.TTL) * time.Second
	if ttl == 0 {
		ttl = 10 * time.Minute // Default TTL: 10 minutes
	}

	return &RedisCache{
		redis: redis,
		ttl:   ttl,
	}
}

// Get retrieves a value from the cache
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := rc.redis.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found in cache")
		}
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in the cache with the default TTL
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return rc.redis.Client.Set(ctx, key, data, rc.ttl).Err()
}

// Delete removes a value from the cache
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.redis.Client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.redis.Client.Close()
}
