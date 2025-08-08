package logging

// Based on https://gist.github.com/panta/2530672ca641d953ae452ecb5ef79d7d

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/alphatoasterous/otlozhka-bot/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration for logging
type Logger struct {
	*zerolog.Logger
}

// Configure sets up the logging framework
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func Configure(config config.ZerologConfiguration) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(mw).With().Timestamp().Logger()

	logger.Info().
		Bool("fileLogging", config.FileLoggingEnabled).
		Bool("jsonLogOutput", config.EncodeLogsAsJson).
		Str("logDirectory", config.Directory).
		Str("fileName", config.Filename).
		Int("maxSizeMB", config.MaxSize).
		Int("maxBackups", config.MaxBackups).
		Int("maxAgeInDays", config.MaxAge).
		Msg("logging configured")

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(config config.ZerologConfiguration) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		fmt.Printf("ERROR: can't create log directory: %v, path: %s\n", err, config.Directory)
		log.Fatal("A fatal error has occured initializing zerolog.")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}

var Log *Logger

func init() {
	// Setting up zerolog logger
	zerologConfig := config.BotConfig.ZerologConfig
	Log = Configure(zerologConfig)
}
