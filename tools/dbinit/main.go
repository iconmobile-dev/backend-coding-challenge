package main

import (
	"flag"
	"os"

	"github.com/iconmobile-dev/go-coding-challenge/lib/bootstrap"
	"github.com/iconmobile-dev/go-coding-challenge/lib/storage"
)

var schemaPath = flag.String("schemaPath", "", "create schema with given sql schema file")

func main() {
	// bootstrap logger and config
	log, cfg := bootstrap.LoggerAndConfig("dbinit", false)

	flag.Parse()

	// open database
	db, err := storage.NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Error(err, "database postgres new")
		os.Exit(1)
	}

	if schemaPath != nil && *schemaPath != "" {
		err = db.Init(*schemaPath)
		if err != nil {
			log.Error(err, "database init")
			os.Exit(1)
		}
	}

	log.Info("database init done")
}
