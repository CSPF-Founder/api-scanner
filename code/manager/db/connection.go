package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(dbURI string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectDBWithRetry(dbURI string, maxRetries int, initialWaitTime time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err = sql.Open("mysql", dbURI)
		if err != nil {
			waitTime := initialWaitTime * time.Duration(attempt)
			time.Sleep(waitTime)
			continue
		}

		err = db.Ping()
		if err == nil {
			return db, nil
		}

		waitTime := initialWaitTime * time.Duration(attempt)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect after %d attempts", maxRetries)
}
