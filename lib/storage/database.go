package storage

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	// importing solely for its side-effects
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// DB sqlx conn wrap.
type DB struct {
	*sqlx.DB
}

// NewDB returns a new connection to the postgres database.
func NewDB(host, port, user, password, database, sslMode string) (*DB, error) {
	creds := fmt.Sprintf("host='%s' user='%s' password='%s' dbname='%s' port='%s' sslmode='%s'",
		host, user, password, database, port, sslMode)

	conn, err := sqlx.Connect("postgres", creds)
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to db")
	}

	return &DB{conn}, nil
}

// Reset tries to truncate existing tables, should NOT be run on production!
func (db *DB) Reset() error {
	sql := "TRUNCATE users RESTART IDENTITY CASCADE"
	if _, err := db.Exec(sql); err != nil {
		return errors.Wrapf(err, "database reset failed: %v", err)
	}
	return nil
}

// Init creates required DB tables, indices, triggers, etc
// if not already existing from file schema path.
func (db *DB) Init(schemaPath string) error {
	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return errors.Wrap(err, "reading db schema.sql")
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return errors.Wrap(err, "exec schema.sql")
	}

	return nil
}

// TableExists returns false if table does not exist
func (db *DB) TableExists(table string) bool {
	found := ""
	sql := "SELECT to_regclass($1);"
	if err := db.Get(&found, sql, table); err != nil {
		return false
	}
	return true
}
