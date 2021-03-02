package bootstrap

import (
	"os"

	"github.com/iconmobile-dev/go-coding-challenge/config"
	"github.com/iconmobile-dev/go-core/logger"
)

// LoggerAndConfig is returning logger & config which are
// required for bootstrapping a service server
// is using CONFIG_FILE env var and if not set uses
// cfgFilePath. If cfgFilePath is not set then it tries to find the config
func LoggerAndConfig(serverName string, test bool) (logger.Logger, config.Config) {
	// init logger
	log := logger.Logger{MinLevel: "verbose"}

	var cfg *config.Config
	var err error
	// load config
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		cfg, err = config.LoadDefaultConfig()
		if err != nil {
			log.Error(err, "config load")
			os.Exit(1)
		}
	} else {
		cfg, err = config.Load(configFile)
		if err != nil {
			log.Error(err, "config load")
			os.Exit(1)
		}
	}

	// set service name
	log.MinLevel = cfg.Logging.MinLevel
	log.TimeFormat = cfg.Logging.TimeFormat
	log.UseColor = cfg.Logging.UseColor
	log.ReportCaller = cfg.Logging.ReportCaller
	cfg.Server.Name = serverName

	if test {
		log.UseColor = false
		cfg.Server.Name += "_test"
	}

	if cfg.Server.Env == "prod" {
		log.UseJSON = true
	}

	return log, *cfg
}
