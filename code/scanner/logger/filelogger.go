package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type FileLogger struct {
	log *zerolog.Logger
}

func GetFileLogger(logFile *os.File) (*FileLogger, error) {
	logHandler := zerolog.New(logFile).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	logHandler.Info().Msg("Logger initialized")
	return &FileLogger{
		log: &logHandler,
	}, nil
}

func (s *FileLogger) Info(msg string) {
	s.log.Info().Msg(msg)
}

func (s *FileLogger) Debug(msg string) {
	s.log.Debug().Msg(msg)
}

func (s *FileLogger) Error(msg string, err error) {
	if err != nil {
		s.log.Error().Err(err).Msg(msg)
	} else {
		s.log.Error().Msg(msg)
	}
}

func (s *FileLogger) Fatal(msg string, err error) {
	if err != nil {
		s.log.Fatal().Err(err).Msg(msg)
	} else {
		s.log.Fatal().Msg(msg)
	}
}

func (s *FileLogger) Warn(msg string, err error) {
	if err != nil {
		s.log.Warn().Err(err).Msg(msg)
	} else {
		s.log.Warn().Msg(msg)
	}
}
