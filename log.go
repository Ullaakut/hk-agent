package main

import (
	"io"
	"strings"

	"github.com/rs/zerolog"
)

// NewZeroLog creates a new zerolog logger
func NewZeroLog(writer io.Writer) *zerolog.Logger {
	zl := zerolog.New(writer).Output(zerolog.ConsoleWriter{Out: writer}).With().Timestamp().Logger()
	return &zl
}

// parseLevel parses a level from string to log level
func parseLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "FATAL":
		return zerolog.FatalLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "WARNING":
		return zerolog.WarnLevel
	case "INFO":
		return zerolog.InfoLevel
	case "DEBUG":
		return zerolog.DebugLevel
	default:
		return zerolog.DebugLevel
	}
}
