package dynamo

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcacheCache struct {
	*Cache
	client *memcache.Client
}

func NewMemcacheCache(client *memcache.Client) *MemcacheCache {
	return &MemcacheCache{
		Cache: &Cache{
			defaultTTL: 5 * time.Minute,
		},
		client: client,
	}
}

func (m *MemcacheCache) Get(key string) (interface{}, error) {
	value, _ := m.Cache.Get(key)
	if value != nil {
		return value, nil
	}

	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}

	value = string(item.Value)
	m.Cache.Set(key, value)
	return value, nil
}

func (m *MemcacheCache) Set(key string, value interface{}) error {
	m.Cache.Set(key, value)
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value.(string)),
		Expiration: int32(m.defaultTTL.Seconds()),
	}
	return m.client.Set(item)
}

func (m *MemcacheCache) Delete(key string) error {
	m.Cache.Delete(key)
	return m.client.Delete(key)
}
