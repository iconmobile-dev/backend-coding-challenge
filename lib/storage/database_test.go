package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TableExists(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
	})

	assert.False(t, db.TableExists("foo"))
	assert.True(t, db.TableExists("users"))
}

func Test_Init(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
	})

	assert.Error(t, db.Init("no/path/to/schema.sql"))
}
