package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewLogProessor(t *testing.T) {
	log := NewZeroLog(bytes.NewBuffer([]byte{}), JSON)
	lp := NewLogProcessor(log, 3, 1024, time.Second, time.Now)

	if lp.topHitsNumber != 3 {
		t.Error("NewLogProcessor doesn't set top hits number properly")
	}
	if lp.trafficThreshold != 1024 {
		t.Error("NewLogProcessor doesn't set traffic threshold properly")
	}
	if lp.log != log {
		t.Error("NewLogProcessor doesn't set logger properly")
	}
	if lp.refreshPeriod != time.Second {
		t.Error("NewLogProcessor doesn't set refresh period properly")
	}
}

func TestAddEntries(t *testing.T) {
	entry1 := &HTTPEntry{
		ClientAddress: "::1",
		Request:       "GET /bestsection/top/one HTTP/1.1",
		UserID:        "frank",
		Identifier:    "user-identifier",
		TimeStr:       "17/May/2054:18:54:34 +0000",
		StatusStr:     "201",
		SizeStr:       "1345",
		Section:       "/bestsection",
		Status:        201,
		Size:          1345,
		Time:          time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
	}
	entry2 := &HTTPEntry{
		ClientAddress: "localhost",
		Request:       "GET /notbestsection/top/two HTTP/1.1",
		UserID:        "frank",
		Identifier:    "user-identifier",
		TimeStr:       "17/May/2054:18:54:34 +0000",
		StatusStr:     "201",
		SizeStr:       "1345",
		Section:       "/notbestsection",
		Status:        201,
		Size:          1345,
		Time:          time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
	}
	entry3 := &HTTPEntry{
		ClientAddress: "::1",
		Request:       "GET /worstsection/top/three HTTP/1.1",
		UserID:        "frank",
		Identifier:    "user-identifier",
		TimeStr:       "17/May/2054:18:54:34 +0000",
		StatusStr:     "201",
		SizeStr:       "1345",
		Section:       "/worstsection",
		Status:        201,
		Size:          1345,
		Time:          time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
	}
	entry4 := &HTTPEntry{
		ClientAddress: "::1",
		Request:       "GET /notevenintop/no/nothing HTTP/1.1",
		UserID:        "frank",
		Identifier:    "user-identifier",
		TimeStr:       "17/May/2054:18:54:34 +0000",
		StatusStr:     "201",
		SizeStr:       "1345",
		Section:       "/notevenintop",
		Status:        201,
		Size:          1345,
		Time:          time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
	}

	b := []byte{}
	buffer := bytes.NewBuffer(b)
	log := NewZeroLog(buffer, JSON)

	lp := &LogProcessor{
		log:              log,
		topHitsNumber:    3,
		trafficThreshold: 1024,
		hits:             make(map[string]int),
		refreshPeriod:    10 * time.Second,
		now:              time.Now,
	}

	entries := []*HTTPEntry{
		entry1, entry1, entry1, entry1, entry1,
		entry2, entry2, entry2, entry2,
		entry3, entry3, entry3,
		entry4,
	}

	lp.Add(entries)

	if !strings.Contains(buffer.String(), `{"level":"info","section":"/bestsection","hits":5,"message":"Top section #1"}`) {
		t.Error(`expected log {"level":"info","section":"/bestsection","hits":5,"message":"Top section #1"}`)
	}
	if !strings.Contains(buffer.String(), `{"level":"info","section":"/notbestsection","hits":4,"message":"Top section #2"}`) {
		t.Error(`expected log {"level":"info","section":"/notbestsection","hits":4,"message":"Top section #2"}`)
	}
	if !strings.Contains(buffer.String(), `{"level":"info","section":"/worstsection","hits":3,"message":"Top section #3"}`) {
		t.Error(`expected log {"level":"info","section":"/worstsection","hits":3,"message":"Top section #3"}`)
	}
	if !strings.Contains(buffer.String(), `{"level":"info","total_entries":13,"recent_entries":13,"message":"Statistics"}`) {
		t.Error(`expected log {"level":"info","total_entries":13,"recent_entries":13,"message":"Statistics"}`)
	}
}

// This test ensures that when the traffic reaches the threshold a message is outputed, and that when it is still the
// case messages are outputed at every iteration of the refresh period. It also makes sure that when the traffic goes
// back below the threshold in the last 2 minutes, the alert stops and a message is outputed to indicate that as well
func TestAlerting(t *testing.T) {
	// December 11, 1241 - look it up
	baseTime := time.Date(1241, time.December, 11, 8, 42, 24, 0, time.UTC)
	entry1 := &HTTPEntry{
		ClientAddress: "::1",
		Request:       "GET /bestsection/top/one HTTP/1.1",
		UserID:        "frank",
		Identifier:    "user-identifier",
		TimeStr:       "17/May/2054:18:54:34 +0000",
		StatusStr:     "201",
		SizeStr:       "999999999",
		Section:       "/bestsection",
		Status:        201,
		Size:          999999999, // 953 MB
		Time:          baseTime,
	}

	b := []byte{}
	buffer := bytes.NewBuffer(b)
	log := NewZeroLog(buffer, JSON)

	lp := &LogProcessor{
		log:              log,
		topHitsNumber:    3,
		trafficThreshold: 1,
		refreshPeriod:    10 * time.Millisecond,
		hits:             make(map[string]int),
		now: func() time.Time {
			// first call: entry is 0s old - alert
			// second call: entry is 1mn30 old - still alert
			// third call: entry is 3mn old - outdated
			baseTime = baseTime.Add(90 * time.Second)
			return baseTime
		},
	}

	entries := []*HTTPEntry{
		entry1,
	}

	lp.Add(entries)
	lp.Add(nil)
	lp.Add(nil)

	if !strings.Contains(buffer.String(), `{"level":"warn","recent_traffic":"953MB","threshold":"1MB","message":"Total traffic over the last 2 minutes exceeds the configured threshold"}`) {
		t.Error(`expected log {"level":"warn","recent_traffic":"953MB","threshold":"1MB","message":"Total traffic over the last 2 minutes exceeds the configured threshold"}`)
	}

	if !strings.Contains(buffer.String(), `{"level":"warn","recent_traffic":"953MB","threshold":"1MB","message":"Total traffic over the last 2 minutes still exceeds the configured threshold"}`) {
		t.Error(`expected log {"level":"warn","recent_traffic":"953MB","threshold":"1MB","message":"Total traffic over the last 2 minutes still exceeds the configured threshold"}`)
	}

	if !strings.Contains(buffer.String(), `{"level":"info","recent_traffic":"0MB","threshold":"1MB","message":"Total traffic over the last 2 minutes is back to normal"}`) {
		t.Error(`expected log {"level":"warn","recent_traffic":"0MB","threshold":"1MB","message":"Total traffic over the last 2 minutes is back to normal"}`)
	}
}
