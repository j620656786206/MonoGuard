package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/monoguard/api/internal/config"
)

// RedisClient wraps the Redis client
type RedisClient struct {
	*redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(cfg *config.RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		
		// Connection pool settings
		PoolSize:        10,
		PoolTimeout:     30 * time.Second,
		IdleTimeout:     5 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute,
		
		// Retry settings
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{client}, nil
}

// HealthCheck checks if Redis is accessible
func (r *RedisClient) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis ping failed: %w", err)
	}

	return nil
}

// SetJSON sets a JSON value in Redis with expiration
func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Set(ctx, key, value, expiration).Err()
}

// GetJSON gets a JSON value from Redis
func (r *RedisClient) GetJSON(ctx context.Context, key string) (string, error) {
	return r.Get(ctx, key).Result()
}

// DeletePattern deletes keys matching a pattern
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.Del(ctx, keys...).Err()
	}

	return nil
}