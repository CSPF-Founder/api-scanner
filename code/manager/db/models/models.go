package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// DBModel is the type for db connection values
// type DBModel struct {
// 	DB *sql.DB
// }

// Models is the wrapper for all models
type Service struct {
	// DB  DBModel
	Job JobService
}

func New(db *sql.DB) Service {
	return Service{
		Job: NewJobService(db),
	}
}

func ParseTimeString(timeString string) (time.Time, error) {
	// Trim any leading or trailing whitespaces
	timeString = strings.TrimSpace(timeString)

	layouts := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05", // Add other formats as needed
	}

	for _, layout := range layouts {
		parsedTime, err := time.Parse(layout, timeString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, fmt.Errorf("Failed to parse time string: %s", timeString)
}
