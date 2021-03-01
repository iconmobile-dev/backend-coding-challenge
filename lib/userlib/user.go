package userlib

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/iconmobile-dev/go-core/errors"
	"github.com/iconmobile-dev/go-core/structs"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/iconmobile-dev/backend-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/backend-coding-challenge/pkg/sqlutil"
)

// User contains the database entry
type User struct {
	ID          int
	Email       string `valid:"required,email"`
	Password    string `json:"-" valid:"required,stringlength(8|99)"`
	Description string
	FirstName   string
	LastName    string
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// Insert sanitizes, validates and inserts a User in database
// Should not be called without prior role check!
func (u *User) Insert(db *storage.DB, cache *storage.Cache) error {
	// removes all leading and trailing white spaces from string fields
	err := u.Sanitize()
	if err != nil {
		return errors.E(err)
	}

	// validate User
	err = u.IsValid(db)
	if err != nil {
		return errors.E(err, errors.Unprocessable)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.E(err, errors.Internal)
	}
	u.Password = string(hashedPassword)

	us, err := ListUsers(UserListParams{
		Filter: UserFilter{
			Email: &sqlutil.StringFilter{
				Is: &u.Email,
			},
		},
	}, db)
	if err != nil {
		return errors.E(err)
	}
	if len(us) != 0 {
		msg := fmt.Sprintf("user with email %v does already exist", u.Email)
		return errors.E(fmt.Errorf(msg), errors.Conflict, msg)
	}

	err = u.insert(db)
	if err != nil {
		return errors.E(err)
	}

	return nil
}

// insert inserts a User in database
// Should not be called without prior data validation!
func (u *User) insert(db *storage.DB) error {
	var returned User
	// insert
	sql := `INSERT INTO users (email, password, firstname, lastname)
			VALUES ($1, $2, $3, $4)
			RETURNING *`
	err := db.Get(&returned, sql, u.Email, u.Password, u.FirstName, u.LastName)
	if err != nil {
		return errors.E(err, errors.Internal)
	}

	*u = returned

	return nil
}

// Update sanitizes, validates and updates User in database
// Should not be called without prior role check!
func (u *User) Update(oldHashedPassword string, oldPassword *string, db *storage.DB, cache *storage.Cache) error {
	// removes all leading and trailing white spaces from string fields
	err := u.Sanitize()
	if err != nil {
		return errors.E(err)
	}

	// validate User
	err = u.IsValid(db)
	if err != nil {
		return errors.E(err, errors.Unprocessable)
	}

	if u.Password != oldHashedPassword && oldPassword != nil {
		err := bcrypt.CompareHashAndPassword([]byte(oldHashedPassword), []byte(*oldPassword))
		if err != nil {
			return errors.E(err, errors.Unprocessable, "OldPassword is incorrect")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("bcrypt.GenerateFromPassword err:", err)
			return errors.E(err, errors.Internal, "Internal server error")
		}

		u.Password = string(hashedPassword)
	}

	err = u.update(db)
	if err != nil {
		return errors.E(err)
	}

	return nil
}

// update updates User in database
// Should not be called without prior data validation!
func (u *User) update(db *storage.DB) error {
	var returned User
	// update
	sql := `UPDATE users
				SET password=$1,
					firstname=$2,
					lastname=$3
				WHERE
					id=$4 RETURNING *`

	err := db.Get(&returned, sql, u.Password, u.FirstName, u.LastName, u.ID)
	if err != nil {
		return errors.E(err, errors.Internal)
	}

	*u = returned

	return nil
}

// IsValid validates an User by checking
// if all field values are valid
func (u User) IsValid(db *storage.DB) error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return errors.E(err, errors.Unprocessable)
	}

	return nil
}

// Sanitize removes all leading and trailing white spaces from string fields
func (u *User) Sanitize() error {
	err := structs.Sanitize(u)
	if err != nil {
		return errors.E(err, errors.Internal)
	}

	return nil
}

// IsCorrectPassword checks if the password is correct
func (u *User) IsCorrectPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.E(err, errors.Unauthorized, "Password is incorrect")
	}

	return nil
}

// UserByID loads User with given ID, returns nil if not found
func UserByID(id int, db *storage.DB) (User, error) {
	u := User{}
	q := `SELECT * FROM users WHERE id=$1 LIMIT 1;`
	if err := db.Get(&u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, errors.E(err, errors.NotFound)
		}
		return u, errors.E(err, errors.Internal)
	}

	return u, nil
}

// UserListParams to list Users
type UserListParams struct {
	Pagination sqlutil.LimitOffsetPagination
	Sort       sqlutil.OneColumnSort
	Filter     UserFilter
}

// UserFilter to filter Users
type UserFilter struct {
	ID            *sqlutil.IntFilter
	Role          *sqlutil.IntFilter
	Email         *sqlutil.StringFilter
	EmailToVerify *sqlutil.StringFilter `db:"email_to_verify"`
	Password      *sqlutil.StringFilter
	FirstName     *sqlutil.StringFilter
	LastName      *sqlutil.StringFilter
	Description   *sqlutil.StringFilter
	ImageURL      *sqlutil.StringFilter `db:"image_url"`
	Language      *sqlutil.StringFilter
	Status        *sqlutil.IntFilter
	Metadata      *sqlutil.StringFilter
	LastLogin     *sqlutil.TimeFilter `db:"last_login"`
	CreatedAt     *sqlutil.TimeFilter `db:"created_at"`
	UpdatedAt     *sqlutil.TimeFilter `db:"updated_at"`
}

// ListUsers returns a list of Users
func ListUsers(params UserListParams, db *storage.DB) ([]User, error) {
	us := []User{}

	q := sqlutil.Select("*").From("users")

	q, err := sqlutil.UseStructFilter(q, "", params.Filter)
	if err != nil {
		return us, errors.E(err)
	}

	q = sqlutil.UseLimitOffsetPagination(q, params.Pagination)

	columnMapping, err := sqlutil.GetColumnMapping(User{})
	if err != nil {
		return us, errors.E(err)
	}

	q, err = sqlutil.UseOneColumnSort(q, params.Sort, columnMapping)
	if err != nil {
		return us, errors.E(err)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return us, errors.E(err, errors.Internal)
	}

	err = db.Select(&us, sql, args...)
	if err != nil {
		return us, errors.E(err, errors.Internal)
	}

	return us, nil
}

func getUserIDs(os []User) []int {
	ids := []int{}

	for _, u := range os {
		ids = append(ids, u.ID)
	}

	return ids
}

// Delete deletes User in database
func (u User) Delete(db *storage.DB) error {
	sql := "DELETE FROM users WHERE id=$1"
	_, err := db.Exec(sql, u.ID)
	if err != nil {
		switch err := err.(type) {
		case *pq.Error:
			switch err.Code.Name() {
			case "foreign_key_violation":
				err := fmt.Errorf("User is still referenced for ID: %d", u.ID)
				return errors.E(err, errors.Unprocessable)
			default:
				break
			}
		default:
			return errors.E(err, errors.Internal)
		}
	}

	return nil
}
