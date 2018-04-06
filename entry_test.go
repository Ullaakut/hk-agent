package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/ullaakut/gonx"
)

// TODO: Test with erroneous entries from gonx for 100% coverage

func TestNewHTTPEntry(t *testing.T) {
	testCases := []struct {
		log string

		expectedClientAddr string
		expectedRequest    string
		expectedUserID     string
		expectedIdentifier string
		expectedTimeStr    string
		expectedStatusStr  string
		expectedSizeStr    string
		expectedSection    string
		expectedStatus     uint64
		expectedSize       uint64
		expectedTime       time.Time
	}{
		{
			log: `179.105.237.248 - - [08/May/2017:08:08:19 +0000] "GET /language/Swedish${IFS}&&echo${IFS}610cker>qt&&tar${IFS}/string.js HTTP/1.0" 404 6407`,

			expectedClientAddr: "179.105.237.248",
			expectedRequest:    "GET /language/Swedish${IFS}&&echo${IFS}610cker>qt&&tar${IFS}/string.js HTTP/1.0",
			expectedUserID:     "-",
			expectedIdentifier: "-",
			expectedTimeStr:    "08/May/2017:08:08:19 +0000",
			expectedStatusStr:  "404",
			expectedSizeStr:    "6407",
			expectedSection:    "/language",
			expectedStatus:     404,
			expectedSize:       6407,
			expectedTime:       time.Date(2017, time.May, 8, 8, 8, 19, 0, time.UTC),
		},
		{
			log: `localhost user-identifier frank [17/May/2054:18:54:34 +0000] "POST / HTTP/1.0" 201 1345`,

			expectedClientAddr: "localhost",
			expectedRequest:    "POST / HTTP/1.0",
			expectedUserID:     "frank",
			expectedIdentifier: "user-identifier",
			expectedTimeStr:    "17/May/2054:18:54:34 +0000",
			expectedStatusStr:  "201",
			expectedSizeStr:    "1345",
			expectedSection:    "/",
			expectedStatus:     201,
			expectedSize:       1345,
			expectedTime:       time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
		},
		{
			log: `::1 user-identifier frank [17/May/2054:18:54:34 +0000] "-" 201 1345`,

			expectedClientAddr: "::1",
			expectedRequest:    "-",
			expectedUserID:     "frank",
			expectedIdentifier: "user-identifier",
			expectedTimeStr:    "17/May/2054:18:54:34 +0000",
			expectedStatusStr:  "201",
			expectedSizeStr:    "1345",
			expectedSection:    "-",
			expectedStatus:     201,
			expectedSize:       1345,
			expectedTime:       time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
		},
		{
			log: `::1 user-identifier frank [17/May/2054:18:54:34 +0000] "OPTIONS * HTTP/1.0" 201 1345`,

			expectedClientAddr: "::1",
			expectedRequest:    "OPTIONS * HTTP/1.0",
			expectedUserID:     "frank",
			expectedIdentifier: "user-identifier",
			expectedTimeStr:    "17/May/2054:18:54:34 +0000",
			expectedStatusStr:  "201",
			expectedSizeStr:    "1345",
			expectedSection:    "*",
			expectedStatus:     201,
			expectedSize:       1345,
			expectedTime:       time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
		},
		{
			log: `::1 user-identifier frank [17/May/2054:18:54:34 +0000] "GET /section/ HTTP/1.0" 201 1345`,

			expectedClientAddr: "::1",
			expectedRequest:    "GET /section/ HTTP/1.0",
			expectedUserID:     "frank",
			expectedIdentifier: "user-identifier",
			expectedTimeStr:    "17/May/2054:18:54:34 +0000",
			expectedStatusStr:  "201",
			expectedSizeStr:    "1345",
			expectedSection:    "/section",
			expectedStatus:     201,
			expectedSize:       1345,
			expectedTime:       time.Date(2054, time.May, 17, 18, 54, 34, 0, time.UTC),
		},
	}
	for _, testCase := range testCases {
		log := NewZeroLog(bytes.NewBuffer([]byte{}), JSON)
		parser := gonx.NewParser(`$client_address $identifier $user_id [$time] "$request" $status $size`)

		entry, err := parser.ParseString(testCase.log)
		if err != nil {
			t.Fatalf("gonx external library failed to parse test log: %s", testCase.log)
		}

		result := NewHTTPEntry(log, entry)

		if result.ClientAddress != testCase.expectedClientAddr {
			t.Errorf("expected client address to be %s, was %s instead", testCase.expectedClientAddr, result.ClientAddress)
		}
		if result.Request != testCase.expectedRequest {
			t.Errorf("expected request to be %s, was %s instead", testCase.expectedRequest, result.Request)
		}
		if result.UserID != testCase.expectedUserID {
			t.Errorf("expected user ID to be %s, was %s instead", testCase.expectedUserID, result.UserID)
		}
		if result.Identifier != testCase.expectedIdentifier {
			t.Errorf("expected identifier to be %s, was %s instead", testCase.expectedIdentifier, result.Identifier)
		}
		if result.TimeStr != testCase.expectedTimeStr {
			t.Errorf("expected time string to be %s, was %s instead", testCase.expectedTimeStr, result.TimeStr)
		}
		if result.StatusStr != testCase.expectedStatusStr {
			t.Errorf("expected StatusStr to be %s, was %s instead", testCase.expectedStatusStr, result.StatusStr)
		}
		if result.SizeStr != testCase.expectedSizeStr {
			t.Errorf("expected SizeStr to be %s, was %s instead", testCase.expectedSizeStr, result.SizeStr)
		}
		if result.Section != testCase.expectedSection {
			t.Errorf("expected Section to be %s, was %s instead", testCase.expectedSection, result.Section)
		}
		if result.Status != testCase.expectedStatus {
			t.Errorf("expected Status to be %d, was %d instead", testCase.expectedStatus, result.Status)
		}
		if result.Size != testCase.expectedSize {
			t.Errorf("expected Size to be %d, was %d instead", testCase.expectedSize, result.Size)
		}
		if !result.Time.Equal(testCase.expectedTime) {
			t.Errorf("expected Time to be %s, was %s instead", testCase.expectedTime, result.Time)
		}
	}
}
