package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	os.Setenv("DSN", "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("REMOTE_WORK_DIR", "/app/data/")
	os.Setenv("LOCAL_TEMP_DIR", "/app/data/temp_uploads/")
	os.Setenv("SCANNER_IMAGE", "scanner-image")
	os.Setenv("PYTHON_SCRIPT_PATH", "/app/bin/reporter")

	os.Setenv("USE_DOTENV", "false")
	config := Load()

	if config.DSN != "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local" {
		t.Errorf("Expected DSN to be 'root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local', got %s", config.DSN)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected LogLevel to be 'info', got %s", config.Logging.Level)
	}

	if config.RemoteWorkDir != "/app/data/" {
		t.Errorf("Expected RemoteWorkDir to be '/app/data/', got %s", config.RemoteWorkDir)
	}

	if config.LocalTempDir != "/app/data/temp_uploads/" {
		t.Errorf("Expected LocalTempDir to be '/app/data/temp_uploads/', got %s", config.LocalTempDir)
	}

	if config.ScannerImage != "scanner-image" {
		t.Errorf("Expected ScannerImage to be 'scanner-image', got %s", config.ScannerImage)
	}

	if config.ReporterBinPath != "/app/bin/reporter" {
		t.Errorf("Expected ReporterBinPath to be '/app/bin/reporter', got %s", config.ReporterBinPath)
	}

}
