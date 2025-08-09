package interfaces

import (
	"context"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Close() error
}
	Increment(ctx context.Context, key string) (int64, error)
	Close() error
	HealthCheck(ctx context.Context) error
	SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetTTL(ctx context.Context, key string) (time.Duration, error)
	Exists(ctx context.Context, key string) (bool, error)
	Flush(ctx context.Context) error
	GetKeys(ctx context.Context, pattern string) ([]string, error)
}
