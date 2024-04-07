package logger

import (
	"os"
	"strings"
	"testing"
)

func TestGetScannerLogger(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "api-scanner-test")
	if err != nil {
		t.Errorf("Error creating temp dir: %v", err)
	}

	defer os.RemoveAll(tmpDir)

	logFilePath := tmpDir + "/test.log"

	scanLogFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Errorf("Error opening log file: %v", err)
	}
	defer scanLogFile.Close()

	l, err := GetFileLogger(scanLogFile)

	if err != nil {
		t.Errorf("Error getting scanner logger: %v", err)
	}
	if l == nil {
		t.Errorf("Logger is nil")
	}

	l.Info("Test log message")

	// Read the file
	// Check if the log message is present

	data, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Errorf("Error reading log file: %v", err)
	}

	if !strings.Contains(string(data), "Test log message") {
		t.Errorf("Log message not found in file: %v", string(data))
	}
}
