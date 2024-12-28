package logger

import (
	"github.com/rs/zerolog"
	"log"

	"os"
)

func ConfigureLogger(environment string) (zerolog.Logger, *os.File) {
	var logger zerolog.Logger
	var file *os.File
	var err error
	const filename = "./logs/app.log"

	// Create logger based on environment
	if environment == "DEV" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	} else {
		// Ensure the logs directory exists
		if err = os.MkdirAll("./logs", os.ModePerm); err != nil {
			log.Fatalf("error creating logs directory: %v", err)
		}

		// Open or create the log file
		if file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			log.Fatalf("error opening/creating log file: %v", err)
		}
		logger = zerolog.New(file).With().Timestamp().Logger()
	}

	return logger, file
}
