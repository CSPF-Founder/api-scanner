package logger

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger(t *testing.T) {
	testCases := []struct {
		name  string
		level string
		want  zerolog.Level
	}{
		{"DebugLevel", "debug", zerolog.DebugLevel},
		{"InfoLevel", "info", zerolog.InfoLevel},
		{"WarnLevel", "warn", zerolog.WarnLevel},
		{"FatalLevel", "fatal", zerolog.FatalLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{Level: tc.level}
			logger, err := NewLogger(config)
			if err != nil {
				t.Fatalf("NewLogger() error = %v, wantErr %v", err, false)
			}

			if logger.Provider.GetLevel() != tc.want {
				t.Errorf("NewLogger() got = %v, want %v", logger.Provider.GetLevel(), tc.want)
			}
		})
	}
}
