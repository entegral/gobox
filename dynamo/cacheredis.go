package dynamo

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	*Cache
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(client *redis.Client, ctx context.Context) CacheInterface {
	return &RedisCache{
		Cache: &Cache{
			defaultTTL: 5 * time.Minute,
		},
		client: client,
		ctx:    ctx,
	}
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	value, _ := r.Cache.Get(key)
	if value != nil {
		return value, nil
	}

	value, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	r.Cache.Set(key, value)
	return value, nil
}

func (r *RedisCache) Set(key string, value interface{}) error {
	r.Cache.Set(key, value)
	return r.client.Set(r.ctx, key, value, r.defaultTTL).Err()
}

func (r *RedisCache) Delete(key string) error {
	r.Cache.Delete(key)
	return r.client.Del(r.ctx, key).Err()
}
