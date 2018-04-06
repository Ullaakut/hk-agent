package main

import (
	"io"
	"strings"

	"github.com/rs/zerolog"
)

// LogMode defines the logging mode.
type LogMode uint8

// LodMode enum to define the different logging modes
const (
	Pretty LogMode = iota
	JSON
)

// NewZeroLog creates a new zerolog logger
func NewZeroLog(writer io.Writer, mode LogMode) *zerolog.Logger {
	var zl zerolog.Logger

	switch mode {
	case JSON:
		zl = zerolog.New(writer)
	case Pretty:
		zl = zerolog.New(writer).Output(zerolog.ConsoleWriter{Out: writer}).With().Timestamp().Logger()
	}

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
