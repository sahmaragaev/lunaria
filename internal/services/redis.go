package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sahmaragaev/lunaria-backend/internal/config"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(cfg *config.RedisConfig) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisService{
		client: client,
	}
}

// BlacklistToken adds a token to the blacklist with expiration
func (r *RedisService) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	return r.client.Set(ctx, key, "revoked", expiration).Err()
}

// IsTokenBlacklisted checks if a token is in the blacklist
func (r *RedisService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// StoreUserSession stores user session data
func (r *RedisService) StoreUserSession(ctx context.Context, userID string, sessionData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", userID)
	return r.client.HSet(ctx, key, sessionData).Err()
}

// GetUserSession retrieves user session data
func (r *RedisService) GetUserSession(ctx context.Context, userID string) (map[string]string, error) {
	key := fmt.Sprintf("session:%s", userID)
	return r.client.HGetAll(ctx, key).Result()
}

// DeleteUserSession removes user session data
func (r *RedisService) DeleteUserSession(ctx context.Context, userID string) error {
	key := fmt.Sprintf("session:%s", userID)
	return r.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (r *RedisService) Close() error {
	return r.client.Close()
}
