package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisUTF8(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, cache.Reset())
	})

	const prefix = "api:test:"
	const expiration = time.Duration(10) * time.Second

	// test if we can use unicode chars
	// in Redis keys and / or values
	key := prefix + "lè Äö 🎉"
	value := "Bonjour lè ÄlØ 🚀"

	err := cache.Set(key, value, expiration).Err()
	assert.Nil(t, err)

	// load that key
	loadedVal, err := cache.Get(key).Result()
	assert.Nil(t, err)
	assert.Equal(t, value, loadedVal)
}
