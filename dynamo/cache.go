package dynamo

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/entegral/gobox/types"
	"github.com/go-redis/redis/v8"
)

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	TTL      time.Duration
}

func generateCacheKey(k types.Keyable) (string, error) {
	partitionKey, sortKey, err := k.Keys(0) // assuming 0 for primary composite key
	if err != nil {
		return "", err
	}

	// Concatenate partitionKey and sortKey to form a unique key
	key := partitionKey + ":" + sortKey

	// Hash the key for standardized cache keys
	h := sha256.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (dbm DBManager) WithCache(dynamoDBClient *dynamodb.Client, cacheConfig *CacheConfig) *DBManager {
	dbm.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port),
		Password: cacheConfig.Password,
		DB:       0, // use default DB
	})
	return &dbm
}

func (d *DBManager) GetWithCache(ctx context.Context, row types.Linkable) (loaded bool, err error) {
	cacheKey, err := generateCacheKey(row)
	if err != nil {
		return false, err
	}

	// Try to get the item from cache
	item, err := d.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && item != "" {
		return true, json.Unmarshal([]byte(item), &row)
	}

	// If not in cache or error occurred, get from DynamoDB
	loaded, err = d.Get(ctx, row)
	if err != nil {
		return loaded, err
	}

	// Put the item in cache for next time
	bstring, _ := json.Marshal(row)
	_ = d.RedisClient.Set(ctx, cacheKey, string(bstring), 0).Err()

	return loaded, nil
}

func (d *DBManager) PutWithCache(ctx context.Context, row types.Linkable) (err error) {
	cacheKey, err := generateCacheKey(row)
	if err != nil {
		return err
	}

	// Put the item in DynamoDB
	err = d.Put(ctx, row)
	if err != nil {
		return err
	}

	// Put the item in cache for next time
	bstring, _ := json.Marshal(row)
	_ = d.RedisClient.Set(ctx, cacheKey, string(bstring), 0).Err()

	return nil
}

func (d *DBManager) DeleteWithCache(ctx context.Context, row types.Linkable) (err error) {
	cacheKey, err := generateCacheKey(row)
	if err != nil {
		return err
	}

	// Delete the item from DynamoDB
	err = d.Delete(ctx, row)
	if err != nil {
		return err
	}

	// Delete the item from cache
	_ = d.RedisClient.Del(ctx, cacheKey).Err()

	return nil
}
