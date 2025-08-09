package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/OmidRasouli/weather-api/config"
	"github.com/OmidRasouli/weather-api/infrastructure/database"
	"github.com/OmidRasouli/weather-api/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// Redis represents a Redis database connection
type Redis struct {
	Client *redis.Client
	TTL    time.Duration
}

// NewRedisConnection creates a new Redis connection from configuration
func NewRedisConnection(cfg config.RedisConfig) (database.RedisClient, error) {
	logger.Infof("Connecting to Redis at %s:%d", cfg.Host, cfg.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, 6379), // Default Redis port is 6379
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Verify connection with longer timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Errorf("Failed to connect to Redis at %s:%d: %v", cfg.Host, cfg.Port, err)
		client.Close() // Clean up failed connection
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	logger.Infof("Successfully connected to Redis: %s", pong)

	ttl := time.Duration(cfg.TTL) * time.Second
	if ttl == 0 {
		ttl = 10 * time.Minute // Default TTL
	}

	return &Redis{
		Client: client,
		TTL:    ttl,
	}, nil
}

// HealthCheck pings the Redis server to check if it's available
func (r *Redis) HealthCheck(ctx context.Context) error {
	if r.Client == nil {
		return fmt.Errorf("redis client is nil")
	}

	// Check if client is closed
	if r.Client.Options().Addr == "" {
		return fmt.Errorf("redis client is closed")
	}

	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		logger.Errorf("Redis health check failed: %v", err)
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

// SetWithTTL sets a key with a custom TTL
func (r *Redis) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return r.Client.Set(ctx, key, data, ttl).Err()
}

// Set sets a key with the default TTL
func (r *Redis) Set(ctx context.Context, key string, value interface{}) error {
	return r.SetWithTTL(ctx, key, value, r.TTL)
}

// Get retrieves a value and unmarshals it to the provided destination
func (r *Redis) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// GetTTL returns the remaining TTL for a key
func (r *Redis) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

// Exists checks if a key exists
func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// Delete removes one or more keys
func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Flush removes all keys
func (r *Redis) Flush(ctx context.Context) error {
	return r.Client.FlushAll(ctx).Err()
}

// GetKeys returns keys matching the pattern
func (r *Redis) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return r.Client.Keys(ctx, pattern).Result()
}

// Expire sets expiration for a key
func (r *Redis) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.Client.Expire(ctx, key, ttl).Err()
}

// Increment atomically increments a key's value
func (r *Redis) Increment(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// Close closes the Redis connection
func (r *Redis) Close() error {
	return r.Client.Close()
}
