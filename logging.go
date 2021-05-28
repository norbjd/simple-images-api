package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	logLevelValue, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevelValue = "info"
	}
	logLevel, err := log.ParseLevel(logLevelValue)
	if err != nil {
		log.Warn("Log level {} unknown, defaulting to info", logLevel)
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)
}
