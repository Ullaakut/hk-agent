package main

import (
	"bufio"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/ullaakut/gonx"
)

func main() {
	// instantiate structured logger
	log := NewZeroLog(os.Stderr, Pretty)

	config := DefaultConfig()
	config.Print(log)

	zerolog.SetGlobalLevel(parseLevel(config.LogLevel))

	// Catch signals
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// read logs in a separate routin
	go readLogs(log, config)

	// Wait for agent to be stopped
	<-sig
	signal.Stop(sig)
	close(sig)
	os.Exit(0)
}

// Reads the logs from the file specified in the configuration
// and process the entries using the configured values
func readLogs(log *zerolog.Logger, config Config) {
	// instanciate parser for Common Log Format
	parser := gonx.NewParser(`$client_address $identifier $user_id [$time] "$request" $status $size`)

	// instantiate log processor
	logProcessor := NewLogProcessor(
		log,
		config.TopHitsNumber,
		config.TrafficThreshold,
		config.RefreshPeriod,
	)

	// open log file
	file, err := os.Open(config.LogFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Could not open logfile")
	}
	defer file.Close()

	for {
		timeEnd := time.Now().Add(config.RefreshPeriod)

		// Recreate scanner at every iteration to ensure that it gets the new data
		scanner := bufio.NewScanner(file)
		entries := []*HTTPEntry{}
		for scanner.Scan() {
			// parse every line of the log file into an HTTP entry
			line := scanner.Text()
			if line != "" {
				entry, err := parser.ParseString(line)
				if err != nil {
					log.Error().Err(err).Msg("Could not parse string")
				} else {
					// convert parsed entry into our own HTTPEntry strucutre
					entries = append(entries, NewHTTPEntry(log, entry))
				}
			}
		}

		// add all parsed entries to logProcessor
		go logProcessor.Add(entries)

		// Sleep for 10 seconds minus the time that this loop took to complete
		// If this loop took more than 10s to complete, sleep will return immediately
		time.Sleep(timeEnd.Sub(time.Now()))
	}
}
