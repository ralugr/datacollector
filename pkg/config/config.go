package config

import (
	"fmt"

	"github.com/ralugr/datacollector/pkg/log"
)

type Config struct {
	AppName  string
	LogLevel log.Level
	// Error may be populated by the ConfigOptions provided to NewApplication
	// to indicate that setup has failed.  NewApplication will return this
	// error if it is set.
	Error error
}

type ConfigOption func(*Config)

func AppName(appName string) ConfigOption {
	return func(cfg *Config) { cfg.AppName = appName }
}

func LogLevel(logLevel log.Level) ConfigOption {
	return func(cfg *Config) {
		if !validInput(logLevel) {
			cfg.Error = fmt.Errorf("Invalid value: %v", logLevel)
			return
		}
		cfg.LogLevel = logLevel
	}
}

func DefaultConfig() Config {
	c := Config{}

	c.AppName = "Test App"
	c.LogLevel = log.DebugLevel
	c.Error = nil

	return c
}

func validInput(level log.Level) bool {
	switch level {
	case log.DebugLevel, log.InfoLevel, log.WarnLevel, log.ErrorLevel:
		return true
	default:
		return false
	}
}
