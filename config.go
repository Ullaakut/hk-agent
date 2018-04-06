package main

import (
	"time"

	"github.com/rs/zerolog"
)

// TODO: Config override using environment

// Config represents the HKAgent configuration
type Config struct {
	LogLevel string
	// File path to the log file that will be read by hk-agent
	LogFilePath string
	// traffic threshold that triggers an alert when traffic from the last 2mns represents more
	// megabytes than this number
	TrafficThreshold uint64
	// number of top hits to display when processing metrics
	TopHitsNumber int
	RefreshPeriod time.Duration
}

// DefaultConfig generates a configuration structure with the default values
func DefaultConfig() Config {
	return Config{
		LogLevel:         "DEBUG",
		LogFilePath:      "logs",
		TrafficThreshold: 4096,
		TopHitsNumber:    3,
		RefreshPeriod:    10 * time.Second,
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
