package sqlutil_test

import (
	"os"
	"testing"
	"time"

	"github.com/iconmobile-dev/go-coding-challenge/config"
	"github.com/iconmobile-dev/go-coding-challenge/lib/bootstrap"
	"github.com/iconmobile-dev/go-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/go-coding-challenge/pkg/sqlutil"
	"github.com/iconmobile-dev/go-core/logger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *storage.DB
var cache *storage.Cache

var log logger.Logger
var cfg config.Config

// SetupLoggerAndConfig sets the global logger and config dependency
// should be called during tests
func SetupLoggerAndConfig(serverName string, test bool) {
	log, cfg = bootstrap.LoggerAndConfig(serverName, test)
}

// initiates log and cfg with default values
func init() {
	SetupLoggerAndConfig("rewardlib", false)
}

func TestMain(m *testing.M) {
	// setup before tests
	var err error

	// bootstrap logger and config
	SetupLoggerAndConfig("rewardlib", true)

	// database
	db, err = storage.NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Error(err, "test database postgres new")
		os.Exit(1)
	}
	log.Verbose("connected to Postgres at", cfg.DB.Host)

	err = db.Reset()
	if err != nil {
		log.Error(err, "database reset")
		os.Exit(1)
	}

	// cache
	cache, err = storage.NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Error(err, "test cache database redis new")
		os.Exit(1)
	}
	log.Verbose("connected to redis at", cfg.DB.Host)

	err = cache.Reset()
	if err != nil {
		log.Error(err, "cache reset")
		os.Exit(1)
	}

	err = createTestTable(db)
	if err != nil {
		log.Error(err, "create test table")
		os.Exit(1)
	}

	// run tests
	code := m.Run()

	// shutdown after tests
	db.Close()
	cache.Close()

	os.Exit(code)
}

type testRow struct {
	ID           int
	TestCase     string     `db:"test_case"`
	IntColumn    *int       `db:"int_column"`
	StringColumn *string    `db:"string_column"`
	BoolColumn   *bool      `db:"bool_column"`
	TimeColumn   *time.Time `db:"time_column"`
}

func createTestTable(db *storage.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS sqlutil_test (
	id serial,
	test_case text NOT NULL,
	int_column int,
	string_column text,
	bool_column boolean,
	time_column timestamp with time zone,
	PRIMARY KEY (id)
)`)
	if err != nil {
		return err
	}

	return nil
}

func mustInsertTestRows(t *testing.T, db *storage.DB, testCase string, rs []testRow) []testRow {
	var returned []testRow
	for _, r := range rs {
		r.TestCase = testCase
		returned = append(returned, mustInsertTestRow(t, db, r))
	}

	return returned
}

func mustInsertTestRow(t *testing.T, db *storage.DB, r testRow) testRow {
	var returned testRow
	err := db.Get(&returned, `
INSERT INTO sqlutil_test (test_case, int_column, string_column, time_column, bool_column)
VALUES ($1, $2, $3, $4, $5)
RETURNING *`, r.TestCase, r.IntColumn, r.StringColumn, r.TimeColumn, r.BoolColumn)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	return returned
}

func resetTestTable(db *storage.DB) error {
	sql := "TRUNCATE sqlutil_test"
	if _, err := db.Exec(sql); err != nil {
		return errors.Wrapf(err, "database reset failed: %v", err)
	}
	return nil
}

func TestGetColumns(t *testing.T) {
	t.Run(`GetColumns`, func(t *testing.T) {
		columns, err := sqlutil.GetColumns(testRowFilter{})
		require.NoError(t, err)

		require.Equal(t, []string{"id", "test_case", "int_column", "string_column", "bool_column", "time_column"}, columns)
	})

	t.Run(`without struct`, func(t *testing.T) {
		_, err := sqlutil.GetColumns(nil)
		require.Error(t, err)
	})
}

func TestGetColumnMapping(t *testing.T) {
	t.Run(`GetColumnMapping`, func(t *testing.T) {
		a := struct {
			ID    int
			OrgID int `db:"org_id"`
		}{}

		expected := map[string]string{
			"ID":     "id",
			"id":     "id",
			"OrgID":  "org_id",
			"org_id": "org_id",
		}

		m, err := sqlutil.GetColumnMapping(a)
		assert.NoError(t, err)
		assert.Equal(t, expected, m)
	})

	t.Run(`without struct`, func(t *testing.T) {
		_, err := sqlutil.GetColumnMapping(nil)
		require.Error(t, err)
	})
}
