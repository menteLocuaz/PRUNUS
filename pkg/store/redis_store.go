package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/prunus/pkg/models"
	"github.com/prunus/pkg/utils/performance"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) models.CacheStore {
	return &RedisStore{client: client}
}

func (s *RedisStore) Get(ctx context.Context, key string, dest interface{}) error {
	defer performance.Trace(ctx, "redis", "Get", performance.RedisThreshold, time.Now())
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	defer performance.Trace(ctx, "redis", "Set", performance.RedisThreshold, time.Now())
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, val, expiration).Err()
}

func (s *RedisStore) Delete(ctx context.Context, key string) error {
	defer performance.Trace(ctx, "redis", "Delete", performance.RedisThreshold, time.Now())
	return s.client.Del(ctx, key).Err()
}

func (s *RedisStore) DeleteByPrefix(ctx context.Context, prefix string) error {
	defer performance.Trace(ctx, "redis", "DeleteByPrefix", performance.RedisThreshold, time.Now())
	var cursor uint64
	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
