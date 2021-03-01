package storage

import (
	"os"
	"testing"
)

var db *DB
var cache *Cache

func TestMain(m *testing.M) {
	// setup before tests
	var err error

	// bootstrap logger and config
	SetupLoggerAndConfig("storage", true)

	// database
	db, err = NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Error(err, "database postgres new")
		os.Exit(1)
	}
	log.Verbose("connected to Postgres at", cfg.DB.Host)

	err = db.Reset()
	if err != nil {
		log.Error(err, "database reset")
		os.Exit(1)
	}

	// cache
	cache, err = NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Error(err, "cache database redis new")
		os.Exit(1)
	}
	log.Verbose("connected to redis at", cfg.DB.Host)

	err = cache.Reset()
	if err != nil {
		log.Error(err, "cache reset")
		os.Exit(1)
	}

	// run tests
	code := m.Run()

	// shutdown after tests
	db.Close()

	os.Exit(code)
}
