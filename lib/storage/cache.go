package storage

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Cache redis client wrap.
type Cache struct {
	*redis.Client
}

// NewCache returns a new connection to the redis database.
func NewCache(host string, port int, password string) (*Cache, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	// check if connection is working by doing a ping
	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to redis")
	}

	return &Cache{client}, nil
}

// Reset flushes the cache
func (cache *Cache) Reset() error {
	_, err := cache.FlushAll().Result()
	if err != nil {
		return errors.Wrap(err, "flushing the cache")
	}
	return nil
}
