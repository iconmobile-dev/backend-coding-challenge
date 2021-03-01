package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loginRoute(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, serverTest.db.Reset())
		assert.NoError(t, serverTest.cache.Reset())
	})

	// myURL := ts.URL + "/users/v1/UserCreate"
}

func Test_logoutRoute(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, serverTest.db.Reset())
		assert.NoError(t, serverTest.cache.Reset())
	})

	// myURL := ts.URL + "/auth/v1/Logout"
}
