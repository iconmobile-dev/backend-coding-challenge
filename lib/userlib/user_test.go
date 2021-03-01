package userlib

import (
	"testing"
	"time"

	"github.com/iconmobile-dev/backend-coding-challenge/pkg/ptrutil"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/sqlutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserInsert(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
		assert.NoError(t, cache.Reset())
	})

	validUser := User{
		Email:     "user0@org.com",
		Password:  "password",
		FirstName: "firstname0",
		LastName:  "lastname0",
	}

	t.Run("insert valid User", func(t *testing.T) {
		user := validUser
		err := user.Insert(db, cache)
		require.NoError(t, err)

		t.Run("assert computed fields", func(t *testing.T) {
			assert.NotZero(t, user.ID)
			assert.WithinDuration(t, time.Now(), user.CreatedAt, 100*time.Millisecond)
			assert.WithinDuration(t, time.Now(), user.UpdatedAt, 100*time.Millisecond)
		})

		t.Run("assert state", func(t *testing.T) {
			// load the user by ID
			loadedUser, err := UserByID(user.ID, db)
			require.NoError(t, err)

			assert.Equal(t, user.ID, loadedUser.ID)
			assert.Equal(t, user.FirstName, loadedUser.FirstName)
			assert.Equal(t, user.CreatedAt, loadedUser.CreatedAt)
			assert.Equal(t, user.UpdatedAt, loadedUser.UpdatedAt)
		})
	})

	t.Run(`insert invalid User with User.Email == ""`, func(t *testing.T) {
		user := validUser
		user.Email = ""
		err := user.Insert(db, cache)
		assert.Error(t, err)
	})

	t.Run(`insert invalid User with User.Email already existing`, func(t *testing.T) {
		user := validUser
		user.Email = "user1@org.com"
		err := user.Insert(db, cache)
		assert.NoError(t, err)

		err = user.Insert(db, cache)
		assert.Error(t, err)
	})
}

func TestUserByID(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
		assert.NoError(t, cache.Reset())
	})

	t.Run("get valid User with ID == insertedUser.ID", func(t *testing.T) {
		insertedUser := User{
			Email:     "maria@example.com",
			Password:  "password",
			FirstName: "Maria",
			LastName:  "Smith",
		}
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		user, err := UserByID(insertedUser.ID, db)
		require.NoError(t, err)

		assert.Equal(t, insertedUser.ID, user.ID)
		assert.Equal(t, insertedUser.FirstName, user.FirstName)
		assert.Equal(t, insertedUser.CreatedAt, user.CreatedAt)
		assert.Equal(t, insertedUser.UpdatedAt, user.UpdatedAt)
	})

	t.Run("fail to get User with db == failingDB", func(t *testing.T) {
		_, err := UserByID(-1, failingDB)
		assert.Error(t, err)
	})

	t.Run("get invalid User with ID == -1", func(t *testing.T) {
		_, err := UserByID(-1, db)
		assert.Error(t, err)
	})

	t.Run("get invalid User with ID == 0", func(t *testing.T) {
		_, err := UserByID(0, db)
		assert.Error(t, err)
	})
}

func TestUserUpdate(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
		assert.NoError(t, cache.Reset())
	})

	validUser := User{
		Email:     "user0@org.com",
		Password:  "password",
		FirstName: "firstname0",
		LastName:  "lastname0",
	}

	t.Run("update valid User", func(t *testing.T) {
		insertedUser := validUser
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		updatedUser := insertedUser
		updatedUser.FirstName = "name_b"
		err = updatedUser.Update(insertedUser.Password, nil, db, cache)
		require.NoError(t, err)

		t.Run("assert computed fields", func(t *testing.T) {
			assert.Equal(t, insertedUser.CreatedAt, updatedUser.CreatedAt)
			assert.True(t, updatedUser.UpdatedAt.After(insertedUser.UpdatedAt))
		})

		loadedUser, err := UserByID(updatedUser.ID, db)
		require.NoError(t, err)

		assert.Equal(t, updatedUser.ID, loadedUser.ID)
		assert.Equal(t, updatedUser.FirstName, loadedUser.FirstName)
		assert.Equal(t, updatedUser.CreatedAt, loadedUser.CreatedAt)
		assert.Equal(t, updatedUser.UpdatedAt, loadedUser.UpdatedAt)
	})

	t.Run("update valid User with optional fields", func(t *testing.T) {
		insertedUser := validUser
		insertedUser.Email = "user1@org.com"
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		updatedUser := insertedUser
		updatedUser.Description = "description_a"
		err = updatedUser.Update(insertedUser.Password, nil, db, cache)
		require.NoError(t, err)

		loadedUser, err := UserByID(updatedUser.ID, db)
		require.NoError(t, err)

		assert.Equal(t, updatedUser.Description, loadedUser.Description)
	})

	t.Run("update valid User with new password", func(t *testing.T) {
		insertedUser := validUser
		insertedUser.Email = "user2@org.com"
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		updatedUser := insertedUser
		updatedUser.Description = "description_a"
		updatedUser.Password = "new_password"
		err = updatedUser.Update(insertedUser.Password, &validUser.Password, db, cache)
		require.NoError(t, err)

		err = updatedUser.IsCorrectPassword("new_password")
		require.NoError(t, err)
	})

	t.Run("update valid User with new password, but incorrect old password", func(t *testing.T) {
		insertedUser := validUser
		insertedUser.Email = "user3@org.com"
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		updatedUser := insertedUser
		updatedUser.Description = "description_a"
		updatedUser.Password = "new_password"
		err = updatedUser.Update(insertedUser.Password, ptrutil.String("incorrect_password"), db, cache)
		require.Error(t, err)
	})

	t.Run("update invalid User", func(t *testing.T) {
		insertedUser := validUser
		insertedUser.Email = "user4@org.com"
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		updatedUser := insertedUser
		updatedUser.Email = ""
		err = updatedUser.Update(insertedUser.Password, nil, db, cache)
		require.Error(t, err)
	})

	t.Run("update valid User with failing DB", func(t *testing.T) {
		insertedUser := validUser
		insertedUser.Email = "user5@org.com"
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		updatedUser := insertedUser
		updatedUser.Description = "description_a"
		err = updatedUser.Update(insertedUser.Password, nil, failingDB, nil)
		require.Error(t, err)
	})
}

func TestUserDelete(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
		assert.NoError(t, cache.Reset())
	})

	validUser := User{
		Email:       "user0@org.com",
		Password:    "password",
		FirstName:   "firstname0",
		LastName:    "lastname0",
		Description: "description0",
	}

	t.Run("delete valid User", func(t *testing.T) {
		insertedUser := validUser
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		err = insertedUser.Delete(db)
		require.NoError(t, err)
	})

	t.Run("delete valid User with db == failingDB", func(t *testing.T) {
		insertedUser := validUser
		err := insertedUser.Insert(db, cache)
		require.NoError(t, err)

		err = insertedUser.Delete(failingDB)
		assert.NotNil(t, err)
	})
}

func TestUserList(t *testing.T) {
	t.Cleanup(func() {
		assert.NoError(t, db.Reset())
		assert.NoError(t, cache.Reset())
	})

	// create Users
	user0 := User{
		Email:       "user0@org.com",
		Password:    "password",
		FirstName:   "firstname0",
		LastName:    "lastname0",
		Description: "description0",
	}
	err := user0.Insert(db, cache)
	require.NoError(t, err)

	user1 := User{
		Email:       "user1@org.com",
		Password:    "password",
		FirstName:   "firstname1",
		LastName:    "lastname1",
		Description: "description1",
	}
	err = user1.Insert(db, cache)
	require.NoError(t, err)

	user2 := User{
		Email:       "user2@org.com",
		Password:    "password",
		FirstName:   "firstname2",
		LastName:    "lastname2",
		Description: "description2",
	}
	err = user2.Insert(db, cache)
	require.NoError(t, err)

	t.Run("fail to lists Users with db == failingDB", func(t *testing.T) {
		params := UserListParams{}
		users, err := ListUsers(params, failingDB)
		require.Error(t, err)
		require.True(t, assert.Equal(t, 0, len(users)))
	})

	t.Run("list Users with UserListParams.Pagination.Limit == 2", func(t *testing.T) {
		params := UserListParams{
			Pagination: sqlutil.LimitOffsetPagination{
				Limit: 2,
			},
		}
		users, err := ListUsers(params, db)
		require.NoError(t, err)
		require.True(t, assert.Equal(t, 2, len(users)))

		assert.Equal(t, user0.ID, users[0].ID)
		assert.Equal(t, user1.ID, users[1].ID)
	})

	t.Run("list Users with UserListParams.Pagination.Limit == 2 and UserListParams.Pagination.Offset == 1", func(t *testing.T) {
		params := UserListParams{
			Pagination: sqlutil.LimitOffsetPagination{
				Limit:  2,
				Offset: 1,
			},
		}
		users, err := ListUsers(params, db)
		require.NoError(t, err)
		require.True(t, assert.Equal(t, 2, len(users)))

		assert.Equal(t, user1.ID, users[0].ID)
		assert.Equal(t, user2.ID, users[1].ID)
	})

	t.Run(`list Users with UserListParams.Pagination.Limit == 2 and UserListParams.Sort.Order == "desc"`, func(t *testing.T) {
		params := UserListParams{
			Pagination: sqlutil.LimitOffsetPagination{
				Limit: 2,
			},
			Sort: sqlutil.OneColumnSort{
				Order: "desc",
			},
		}
		users, err := ListUsers(params, db)
		require.NoError(t, err)
		require.True(t, assert.Equal(t, 2, len(users)))

		assert.Equal(t, user2.ID, users[0].ID)
		assert.Equal(t, user1.ID, users[1].ID)
	})

	t.Run(`list Users with UserListParams.Sort.Column == "not_allowed"`, func(t *testing.T) {
		params := UserListParams{
			Sort: sqlutil.OneColumnSort{
				Order:  "desc",
				Column: "not_allowed",
			},
		}
		_, err := ListUsers(params, db)
		require.Error(t, err)
	})

	t.Run(`list Users with UserListParams.OrgID == orgB.ID`, func(t *testing.T) {
		params := UserListParams{
			Filter: UserFilter{
				FirstName: &sqlutil.StringFilter{
					Is: &user0.FirstName,
				},
			},
		}
		users, err := ListUsers(params, db)
		require.NoError(t, err)
		require.True(t, assert.Equal(t, 1, len(users)))

		assert.Equal(t, user0.FirstName, users[0].FirstName)
	})

	t.Run(`list Users return empty, not nil`, func(t *testing.T) {
		filter := UserListParams{
			Pagination: sqlutil.LimitOffsetPagination{
				Offset: 1000,
			},
		}
		emptyReturn, err := ListUsers(filter, db)
		require.NoError(t, err)
		assert.Equal(t, []User{}, emptyReturn)
	})
}
