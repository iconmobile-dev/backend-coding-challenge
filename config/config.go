package config

import (
	_ "embed"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config used globally
type Config struct {
	Server  Server
	Logging Logging
	DB      Database `toml:"database"`
	Redis   Redis
	Crypto  Crypto
}

// Server configuration
type Server struct {
	Env              string
	Name             string
	AssetUploadMaxMB int
	PortGateway      int
	PortEngagement   int
	PortImager       int
	PortProduct      int
}

// Logging configuration
type Logging struct {
	MinLevel     string
	TimeFormat   string
	UseColor     bool
	ReportCaller bool
}

// Database configuration
type Database struct {
	Host     string
	Name     string
	User     string
	Password string
	Port     string
	SSLMode  string
}

// Redis configuration
type Redis struct {
	Host     string
	Port     int
	Password string
}

// Crypto contains encryption keys
type Crypto struct {
	TokenValuePassword string
}

//go:embed config_dev.toml
var configDev string

func LoadDefaultConfig() (*Config, error) {
	var c Config
	_, err := toml.Decode(configDev, &c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode default config")
	}

	return &c, nil
}

// Load reads info from TOML file at relative path
func Load(filename string) (*Config, error) {
	configFile := filename
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "config file is missing: %s", configFile)
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return nil, errors.Wrapf(err, "parsing config file: %s", configFile)
	}

	/*
		// just do the secret value injection on production
		if config.Server.Env == "prod" {
			config = addSecrets(config)
		}*/

	return &config, nil
}
