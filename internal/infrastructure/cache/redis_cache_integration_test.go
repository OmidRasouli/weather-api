package cache_test

import (
	"context"
	"os"
	"testing"

	"github.com/OmidRasouli/weather-api/config"
	"github.com/OmidRasouli/weather-api/infrastructure/database/cache"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_SetAndGet(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		t.Skip("REDIS_ADDR not set, skipping integration test")
	}
	client, err := cache.NewRedisConnection(config.RedisConfig{
		Host: addr,
		Port: 6379,
		DB:   0,
		TTL:  600,
	})
	assert.NoError(t, err)
	defer client.Close()

	ctx := context.Background()
	key := "test-key"
	value := "test-value"

	err = client.Set(ctx, key, value)
	assert.NoError(t, err)

	var got string
	err = client.Get(ctx, key, &got)
	assert.NoError(t, err)
	assert.Equal(t, value, got)

	_ = client.Delete(ctx, key)
}
