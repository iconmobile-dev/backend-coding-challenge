// Package main contains the binary
//
package main

import (
	"fmt"
	"os"

	"github.com/iconmobile-dev/backend-coding-challenge/lib/bootstrap"
	"github.com/iconmobile-dev/backend-coding-challenge/lib/storage"
	"github.com/iconmobile-dev/backend-coding-challenge/services/user"
)

func main() {
	// bootstrap logger and config
	log, cfg := bootstrap.LoggerAndConfig("user", false)

	// open database
	db, err := storage.NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Error(err, "database postgres new")
		os.Exit(1)
	}

	// open cache
	cache, err := storage.NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Error(err, "cache database redis new")
		os.Exit(1)
	}

	// init service
	s := user.New(db, cache)

	log.Info("Starting", cfg.Server.Name, "on", cfg.Server.Env, "using port", cfg.Server.PortEngagement)

	bind := fmt.Sprintf(":%v", cfg.Server.PortEngagement)
	srv := bootstrap.Server(s, bind)

	err = srv.ListenAndServe()
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}
}
