package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog"
)

type hit struct {
	key   string
	value int
}

// LogProcessor is a  structure that contains all previous HTTP logs and processes
// them to detect high traffic and rank top hits for example
type LogProcessor struct {
	log *zerolog.Logger

	// Configuration
	trafficThreshold uint64
	topHitsNumber    int

	// entries from the last 1mn50s, used for calculating the recent traffic
	recent []*HTTPEntry
	// previous state of the hits (avoid recalculating everything at every iteration)
	hits map[string]int
	// current state of the traffic alert
	trafficAlert bool
	// total number of HTTP entries
	totalEntries int
	// entries in the last 2mn
	recentEntries int
}

// NewLogProcessor returns an instance of LogProcessor using the given configuration values
func NewLogProcessor(log *zerolog.Logger, topHitsNumber int, trafficThreshold uint64) *LogProcessor {
	return &LogProcessor{
		log:              log,
		topHitsNumber:    topHitsNumber,
		trafficThreshold: trafficThreshold,
		hits:             make(map[string]int),
	}
}

func (lp *LogProcessor) add(entries []*HTTPEntry) {
	sortedData := make(map[string][]*HTTPEntry)

	// sort hits by section
	for _, entry := range entries {
		sortedData[entry.Section] = append(sortedData[entry.Section], entry)
	}

	lp.totalEntries += len(entries)

	lp.checkRecentTraffic(entries)
	lp.processMetrics(sortedData)
}

// Processes the metrics from the current state of the log processor and the new entries
func (lp *LogProcessor) processMetrics(sortedData map[string][]*HTTPEntry) {
	var newHits []hit
	var topHits []hit

	// create hit structures for newly received data
	for key, value := range sortedData {
		newHits = append(newHits, hit{key: key, value: len(value)})
	}

	// Add new hits to previous hit state
	for _, section := range newHits {
		lp.hits[section.key] += section.value
	}

	// calculate top hits for
	for key, value := range lp.hits {
		topHits = append(topHits, hit{key: key, value: value})
	}

	// sort sections by number of hits
	sort.Slice(topHits, func(i, j int) bool {
		return topHits[i].value > topHits[j].value
	})

	for idx, section := range topHits {
		if idx == lp.topHitsNumber {
			break
		}
		// print number of hits and position for each section until
		// topHitsNumber is reached
		lp.log.Info().
			Str("section", section.key).
			Int("hits", section.value).
			Msgf("Top section #%d", idx+1)
	}
	lp.log.Info().
		Int("total_entries", lp.totalEntries).
		Int("recent_entries", lp.recentEntries).
		Msg("Statistics")
}

// Prints a warning if the recent traffic is above the configured threshold, and as long as it is the case
// Prints an information message when the traffic goes back below the threshold.
func (lp *LogProcessor) checkRecentTraffic(entries []*HTTPEntry) {
	lp.recentEntries = 0
	recentTraffic := uint64(0)

	// Process data previously set as recent entries (last 1mn50)
	for _, entry := range lp.recent {
		lp.recentEntries++
		recentTraffic += entry.Size
	}

	recentLimit := time.Now().Add(-110 * time.Second) // 1mn50s ago
	// Process new entries (last 10s)
	// and store new entries that are within last 1mn50 into recent entries
	for _, entry := range entries {
		if entry.Time.After(recentLimit) {
			lp.recentEntries++
			recentTraffic += entry.Size
			lp.recent = append(lp.recent, entry)
		}
	}

	// convert bytes to MB
	recentTrafficMB := recentTraffic / (1024 * 1024)
	if recentTrafficMB < lp.trafficThreshold {
		if lp.trafficAlert {
			lp.log.Info().
				Str("recent_traffic", fmt.Sprint(recentTrafficMB, "MB")).
				Str("threshold", fmt.Sprint(lp.trafficThreshold, "MB")).
				Msg("Total traffic over the last 2 minutes is back to normal")
			lp.trafficAlert = false
		}
	} else {
		if lp.trafficAlert {
			lp.log.Warn().
				Str("recent_traffic", fmt.Sprint(recentTrafficMB, "MB")).
				Str("threshold", fmt.Sprint(lp.trafficThreshold, "MB")).
				Msg("Total traffic over the last 2 minutes still exceeds the configured threshold")
		} else {
			lp.log.Warn().
				Str("recent_traffic", fmt.Sprint(recentTrafficMB, "MB")).
				Str("threshold", fmt.Sprint(lp.trafficThreshold, "MB")).
				Msg("Total traffic over the last 2 minutes exceeds the configured threshold")
			lp.trafficAlert = true
		}
	}
}
