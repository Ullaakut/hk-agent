package main

import (
	"time"

	"github.com/rs/zerolog"
)

// TODO: Config override using environment

// Config represents the HKAgent configuration
type Config struct {
	// log level used by the logger
	LogLevel string

	// file path to the log file that will be read by hk-agent
	LogFilePath string

	// traffic threshold that triggers an alert when traffic from the last 2mns represents more
	// megabytes than this number
	TrafficThreshold uint64

	// number of top hits to display when processing metrics
	TopHitsNumber int

	// period after which the agent should fetch new logs and display new metrics/alerts
	RefreshPeriod time.Duration
}

// DefaultConfig generates a configuration structure with the default values
func DefaultConfig() Config {
	return Config{
		LogLevel:         "DEBUG",
		LogFilePath:      "logs",
		TrafficThreshold: 1,
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
		Uint64("traffic_threshold", c.TrafficThreshold).
		Int("top_hits_number", c.TopHitsNumber).
		Msg("Configuration")
}
