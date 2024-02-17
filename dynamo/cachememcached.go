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

func (m *MemcacheCache) CheckCache(key string) (interface{}, error) {
	value, _ := m.Cache.CheckCache(key)
	if value != nil {
		return value, nil
	}

	item, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}

	value = string(item.Value)
	m.Cache.SetCache(key, value)
	return value, nil
}

func (m *MemcacheCache) Set(key string, value interface{}) error {
	m.Cache.SetCache(key, value)
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(value.(string)),
		Expiration: int32(m.defaultTTL.Seconds()),
	}
	return m.client.Set(item)
}

func (m *MemcacheCache) Delete(key string) error {
	m.Cache.DeleteCache(key)
	return m.client.Delete(key)
}
