package main

import (
	"time"

	"github.com/rs/zerolog"
)

// TODO: Config override using environment

// Config represents the HKAgent configuration
type Config struct {
	LogLevel      string
	LogFilePath   string
	RefreshPeriod time.Duration
}

// DefaultConfig generates a configuration structure with the default values
func DefaultConfig() Config {
	return Config{
		LogLevel:      "DEBUG",
		LogFilePath:   "logs",
		RefreshPeriod: 10 * time.Second,
	}
}

// Print prints the current configuration
func (c Config) Print(log *zerolog.Logger) {
	log.Debug().
		Str("log_level", c.LogLevel).
		Str("log_file_path", c.LogFilePath).
		Dur("refresh_period", c.RefreshPeriod).
		Msg("Configuration")
}
