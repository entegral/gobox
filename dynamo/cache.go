package dynamo

import (
	"sync"
	"time"
)

type CacheInterface interface {
	SetTTL(time.Time) *UnixTime
	GetTTL() *UnixTime
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
}

type Cache struct {
	TTL        *UnixTime `dynamodbav:"ttl,omitempty" json:"ttl,omitempty"`
	mem        sync.Map
	defaultTTL time.Duration
}

func NewCache(defaultTTL time.Duration) CacheInterface {
	return &Cache{
		defaultTTL: defaultTTL,
		mem:        sync.Map{},
	}
}

func (c *Cache) SetTTL(t time.Time) *UnixTime {
	c.TTL = &UnixTime{t}
	return c.TTL
}

func (c *Cache) GetTTL() *UnixTime {
	return c.TTL
}

func (c *Cache) Get(key string) (interface{}, error) {
	value, _ := c.mem.Load(key)
	return value, nil
}

func (c *Cache) Set(key string, value interface{}) error {
	c.mem.Store(key, value)
	return nil
}

func (c *Cache) Delete(key string) error {
	c.mem.Delete(key)
	return nil
}
