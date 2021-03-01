package user

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_userCreateRoute(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, serverTest.db.Reset())
		assert.NoError(t, serverTest.cache.Reset())
	})

	createURL := ts.URL + "/users/v1/UserCreate"
	failingDBCreateURL := failingDBTs.URL + "/users/v1/UserCreate"

	validCreateReq := userCreateRequest{
		Email:    "user_create0@example.com",
		Password: "password",
	}

	t.Run("valid CreateRequest", func(t *testing.T) {
		t.Parallel()

		createReq := validCreateReq
		createReq.Email = "user_create1@example.com"
		resp := mustPostRequest(t, createURL, createReq, 200)

		var createRsp userResponse
		mustLoadFromResponse(t, resp, &createRsp)

		// assert computed fields
		assert.NotZero(t, createRsp.User.ID)
		assert.WithinDuration(t, time.Now(), createRsp.User.CreatedAt, 1*time.Second)
		assert.WithinDuration(t, time.Now(), createRsp.User.UpdatedAt, 1*time.Second)

		// assert request fields
		assert.Equal(t, "", createRsp.User.Email)
		assert.Equal(t, "", createRsp.User.Password)
	})

	t.Run("invalid CreateRequest with Email already taken", func(t *testing.T) {
		t.Parallel()

		createReq := validCreateReq
		createReq.Email = "user_create2@example.com"
		_ = mustPostRequest(t, createURL, createReq, 200)

		createReq = validCreateReq
		createReq.Email = "user_create2@example.com"
		_ = mustPostRequest(t, createURL, createReq, 409)
	})

	t.Run("invalid CreateRequest with invalid Email", func(t *testing.T) {
		t.Parallel()

		createReq := validCreateReq
		createReq.Email = "user_create@example"
		_ = mustPostRequest(t, createURL, createReq, 422)
	})

	t.Run("invalid CreateRequest with invalid json", func(t *testing.T) {
		t.Parallel()

		_ = mustPostRequest(t, createURL, "text", 400)
	})

	t.Run("valid CreateRequest with failingDB", func(t *testing.T) {
		createReq := validCreateReq
		_ = mustPostRequest(t, failingDBCreateURL, createReq, 500)
	})
}

func Test_userGetRoute(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, serverTest.db.Reset())
		assert.NoError(t, serverTest.cache.Reset())
	})

	validCreateReq := userCreateRequest{
		Email:    "user_get0@example.com",
		Password: "password",
	}

	createURL := ts.URL + "/users/v1/UserCreate"
	getURL := ts.URL + "/users/v1/UserGet"
	//failingDBGetURL := failingDBTs.URL + "/users/v1/UserGet"

	t.Run("valid GetRequest", func(t *testing.T) {
		t.Parallel()

		createReq := validCreateReq
		createReq.Email = "user_get1@example.com"
		resp := mustPostRequest(t, createURL, createReq, 200)
		var createRsp userResponse
		mustLoadFromResponse(t, resp, &createRsp)

		getReq := userGetRequest{
			ID: createRsp.User.ID,
		}
		resp = mustPostRequest(t, getURL, getReq, 200)
		var getRsp userResponse
		mustLoadFromResponse(t, resp, &getRsp)

		assert.Equal(t, "", getRsp.User.Password)
	})

	t.Run("invalid GetRequest with Id == 0", func(t *testing.T) {
		t.Parallel()

		getReq := userGetRequest{
			ID: 0,
		}
		_ = mustPostRequest(t, getURL, getReq, 404)
	})

	t.Run("invalid GetRequest with invalid json", func(t *testing.T) {
		t.Parallel()

		_ = mustPostRequest(t, getURL, "text", 400)
	})
}

func Test_userDeleteRoute(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, serverTest.db.Reset())
		assert.NoError(t, serverTest.cache.Reset())
	})

	createURL := ts.URL + "/users/v1/UserCreate"
	deleteURL := ts.URL + "/users/v1/UserDelete"

	validCreateReq := userCreateRequest{
		Email:    "user_delete0@example.com",
		Password: "password",
	}

	t.Run("valid DeleteRequest", func(t *testing.T) {
		t.Parallel()

		createReq := validCreateReq
		createReq.Email = "user_delete1@example.com"
		resp := mustPostRequest(t, createURL, createReq, 200)

		var createRsp userResponse
		mustLoadFromResponse(t, resp, &createRsp)

		deleteReq := userDeleteRequest{}
		deleteReq.ID = createRsp.User.ID
		_ = mustPostRequest(t, deleteURL, deleteReq, 200)
	})

	t.Run("invalid DeleteRequest with .ID == 0", func(t *testing.T) {
		t.Parallel()

		deleteReq := userDeleteRequest{}
		deleteReq.ID = 0
		_ = mustPostRequest(t, deleteURL, deleteReq, 404)
	})

	t.Run("invalid DeleteRequest with invalid json", func(t *testing.T) {
		t.Parallel()

		_ = mustPostRequest(t, deleteURL, "text", 400)
	})
}
