package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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

	Section string
	Status  uint64
	Size    uint64
	Time    time.Time
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

	h.Section, err = parseSection(log, h.Request)
	if err != nil {
		log.Warn().Err(err).Msg("could not parse section")
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
		Str("section", httpEntry.Section).
		Uint64("status", httpEntry.Status).
		Uint64("size", httpEntry.Size).
		Time("timestamp", httpEntry.Time).
		Msg("Request received")

	return httpEntry
}

func parseSection(log *zerolog.Logger, request string) (string, error) {
	// if the request isn't provided in the logs, no section can be parsed
	if request == "-" {
		return "-", nil
	}

	// Find beginning of section string
	sectionPos := strings.Index(request, " ") + 1
	// Find end of section string
	endOfSectionPos := strings.Index(request[sectionPos:], " ")
	// Find subsection position if it exists
	subsectionPos := strings.Index(request[sectionPos+1:], "/") + 1

	// Request string does not contain the format "METHOD /route/ ..."
	if sectionPos == -1 || endOfSectionPos == -1 {
		return "", errors.New("invalid request format")
	}

	// Remove / from section name
	if request[sectionPos+endOfSectionPos] == '/' {
		endOfSectionPos++
	}

	// if there is no subsection, return the whole route
	if endOfSectionPos < subsectionPos {
		return request[sectionPos : endOfSectionPos+sectionPos], nil
	}
	// otherwise, return only the section
	return request[sectionPos : subsectionPos+sectionPos], nil
}
