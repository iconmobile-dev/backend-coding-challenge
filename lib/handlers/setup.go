package handlers

import (
	"github.com/iconmobile-dev/go-coding-challenge/config"
	"github.com/iconmobile-dev/go-coding-challenge/lib/bootstrap"
	"github.com/iconmobile-dev/go-core/logger"
)

var log logger.Logger
var cfg config.Config

// SetupLoggerAndConfig sets the global logger and config dependency
// should be called during tests
func SetupLoggerAndConfig(serverName string, test bool) {
	log, cfg = bootstrap.LoggerAndConfig(serverName, test)
}

// initiates log and cfg with default values
func init() {
	SetupLoggerAndConfig("handlers", false)
}
