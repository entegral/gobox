package dynamo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestGetWithCache(t *testing.T) {
	ctx := context.Background()

	t.Run("item in cache", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		manager := &DBManager{RedisClient: db}

		row := &User{Email: TEST_EMAIL, Name: TEST_USER_NAME, Age: TEST_USER_AGE}
		cacheKey, _ := generateCacheKey(row)
		bstring, _ := json.Marshal(row)

		mock.ExpectGet(cacheKey).SetVal(string(bstring))

		loaded, err := manager.GetWithCache(ctx, row)

		assert.NoError(t, err)
		assert.True(t, loaded)
		assert.Equal(t, TEST_EMAIL, row.Email)
		assert.Equal(t, TEST_USER_NAME, row.Name)
		assert.Equal(t, TEST_USER_AGE, row.Age)
	})

	t.Run("item not in cache", func(t *testing.T) {
		db, mock := redismock.NewClientMock()
		manager := &DBManager{RedisClient: db}

		row := CreateUser(TEST_EMAIL)
		row.Name = TEST_USER_NAME
		row.Age = TEST_USER_AGE
		cacheKey, _ := generateCacheKey(row)
		bstring, _ := json.Marshal(row)

		mock.ExpectGet(cacheKey).RedisNil()
		mock.ExpectSet(cacheKey, string(bstring), 0).SetVal("")

		loaded, err := manager.GetWithCache(ctx, row)

		assert.NoError(t, err)
		assert.True(t, loaded)
		assert.Equal(t, TEST_EMAIL, row.Email)
		assert.Equal(t, TEST_USER_NAME, row.Name)
		assert.Equal(t, TEST_USER_AGE, row.Age)
	})

	t.Run("error generating cache key", func(t *testing.T) {
		db, _ := redismock.NewClientMock()
		manager := &DBManager{RedisClient: db}

		row := &User{Email: "", Name: TEST_USER_NAME, Age: TEST_USER_AGE}

		loaded, err := manager.GetWithCache(ctx, row)

		assert.Error(t, err)
		assert.False(t, loaded)
	})
}
