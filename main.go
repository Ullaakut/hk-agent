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
	log := NewZeroLog(os.Stderr)

	config := DefaultConfig()
	config.Print(log)

	zerolog.SetGlobalLevel(parseLevel(config.LogLevel))

	// Catch signals
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	file, err := os.Open(config.LogFilePath)
	if err != nil {
		log.Error().Err(err).Msg("Could not open logfile")
	}
	defer file.Close()

	parser := gonx.NewParser(`$client_address $identifier $user_id [$time] "$request" $status $size`)

	go func() {
		for {
			entries := []*HTTPEntry{}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				entry, err := parser.ParseString(scanner.Text())
				if err != nil {
					log.Error().Err(err).Msg("Could not parse string")
				} else {
					entries = append(entries, NewHTTPEntry(log, entry))
				}
			}
			// TODO: Display top hits here
			// TODO: Check traffic on last two minutes
			// TODO: Interface system that keeps alarm information visible at all times
			time.Sleep(config.RefreshPeriod)
		}
	}()

	// Wait for agent to be stopped
	<-sig
	signal.Stop(sig)
	close(sig)
	os.Exit(0)
}
