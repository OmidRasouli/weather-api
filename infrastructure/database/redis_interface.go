package database

import (
	"context"
	"time"
)

// RedisClient defines the interface for Redis operations
type RedisClient interface {
	HealthCheck(ctx context.Context) error
	SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string, dest interface{}) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, keys ...string) error
	Flush(ctx context.Context) error
	GetKeys(ctx context.Context, pattern string) ([]string, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	Increment(ctx context.Context, key string) (int64, error)
	Close() error
}
