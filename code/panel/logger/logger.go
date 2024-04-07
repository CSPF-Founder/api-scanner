package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the main logger that is abstracted in this package.
type Logger struct {
	Provider zerolog.Logger
}

var fallBackLogger *Logger

// Config represents configuration details for logging.
type Config struct {
	Filename string `json:"filename"`
	Level    string `json:"level"`
}

func init() {
	// Initialize the fallback logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.UTC)
	}
	provider := zerolog.New(os.Stdout).With().Timestamp().Logger()

	fallBackLogger = &Logger{Provider: provider}
}

// GetFallBackLogger returns the fallback logger.
func GetFallBackLogger() *Logger {
	return fallBackLogger
}

func NewLogger(config *Config) (*Logger, error) {
	logLevel := getLogLevelFromString(config.Level)
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.UTC)
	}

	provider := zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Logger()
	return &Logger{Provider: provider}, nil
}

func getLogLevelFromString(logLevel string) zerolog.Level {
	switch logLevel {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *Logger) Info(msg string) {
	l.Provider.Info().Msg(msg)
}

func (l *Logger) Debug(msg string) {
	l.Provider.Debug().Msg(msg)
}

func (l *Logger) Error(msg string, err error) {
	if err != nil {
		l.Provider.Error().Err(err).Msg(msg)
	} else {
		l.Provider.Error().Msg(msg)
	}
}

func (l *Logger) Fatal(msg string, err error) {
	if err != nil {
		l.Provider.Fatal().Err(err).Msg(msg)
	} else {
		l.Provider.Fatal().Msg(msg)
	}
}

func (l *Logger) Warn(msg string, err error) {
	if err != nil {
		l.Provider.Warn().Err(err).Msg(msg)
	} else {
		l.Provider.Warn().Msg(msg)
	}
}
