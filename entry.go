package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/ullaakut/gonx"
)

// HTTPEntry represents an entry in an HTTP log file
type HTTPEntry struct {
	ClientAddress string `json:"client_address"`
	Identifier    string `json:"identifier"`
	UserID        string `json:"user_id"`
	Request       string `json:"request"`
	TimeStr       string `json:"time"`
	StatusStr     string `json:"status"`
	SizeStr       string `json:"size"`

	Status uint64
	Size   uint64
	Time   time.Time
}

func (h *HTTPEntry) parseStrings(log *zerolog.Logger) {
	var err error

	h.Time, err = time.Parse(`02/Jan/2006:15:04:05 -0700`, h.TimeStr)
	if err != nil {
		log.Warn().Err(err).Msg("could not parse time")
	}

	h.Status, err = strconv.ParseUint(h.StatusStr, 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("could not parse time")
	}

	h.Size, err = strconv.ParseUint(h.SizeStr, 10, 64)
	if err != nil {
		log.Warn().Err(err).Msg("could not parse time")
	}
}

// NewHTTPEntry instanciates a new HTTPEntry from a gonx.Entry
func NewHTTPEntry(log *zerolog.Logger, entry *gonx.Entry) *HTTPEntry {
	jsonStr, err := entry.ToJSON()
	if err != nil {
		log.Warn().Err(err).Msg("could not parse log entry")
	}

	httpEntry := &HTTPEntry{}
	err = json.Unmarshal(jsonStr, httpEntry)
	if err != nil {
		log.Warn().Err(err).Msg("could not unmarshal log entry into HTTP entry")
	}

	httpEntry.parseStrings(log)

	log.Info().
		Str("client_address", httpEntry.ClientAddress).
		Str("identifier", httpEntry.Identifier).
		Str("user_id", httpEntry.UserID).
		Str("request", httpEntry.Request).
		Uint64("status", httpEntry.Status).
		Uint64("size", httpEntry.Size).
		Time("timestamp", httpEntry.Time).
		Msg("Request received")

	return httpEntry
}
