package config

import (
	"os"
	"strconv"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("PRODUCT_TITLE", "API Scanner")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0:8080")
	os.Setenv("DATABASE_URI", "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local")
	os.Setenv("DBMS_TYPE", "mysql")
	os.Setenv("COPYRIGHT_FOOTER_COMPANY", "Cyber Security & Privacy Foundation")
	os.Setenv("WORK_DIR", "/app/data/")
	os.Setenv("TEMP_UPLOADS_DIR", "/app/data/temp_uploads/")
	os.Setenv("MIGRATIONS_PREFIX", "db")
	os.Setenv("LOG_FILENAME", "logfile")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("USE_TLS", "true")
	os.Setenv("CERT_PATH", "/app/certs/panel.crt")
	os.Setenv("KEY_PATH", "/app/certs/panel.key")

	os.Setenv("USE_DOTENV", "false")
	config, err := LoadConfig()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if config.ProductTitle != "API Scanner" {
		t.Errorf("Expected ProductTitle to be 'API Scanner', got %s", config.ProductTitle)
	}

	if config.ServerConf.ServerAddress != "0.0.0.0:8080" {
		t.Errorf("Expected ServerAddress to be '0.0.0.0:8080', got %s", config.ServerConf.ServerAddress)
	}

	if config.DatabaseURI != "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local" {
		t.Errorf("Expected DatabaseURI to be 'root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local', got %s", config.DatabaseURI)
	}

	if config.DBMSType != "mysql" {
		t.Errorf("Expected DBMSType to be 'mysql', got %s", config.DBMSType)
	}

	if config.CopyrightFooterCompany != "Cyber Security & Privacy Foundation" {
		t.Errorf("Expected CopyrightFooterCompany to be 'Cyber Security & Privacy Foundation', got %s", config.CopyrightFooterCompany)
	}

	if config.WorkDir != "/app/data/" {
		t.Errorf("Expected WorkDir to be '/app/data/', got %s", config.WorkDir)
	}

	if config.TempUploadsDir != "/app/data/temp_uploads/" {
		t.Errorf("Expected TempUploadsDir to be '/app/data/temp_uploads/', got %s", config.TempUploadsDir)
	}

	if config.Logging.Filename != "logfile" {
		t.Errorf("Expected LogFilename to be 'logfile', got %s", config.Logging.Filename)
	}

	if config.Logging.Level != "info" {
		t.Errorf("Expected LogLevel to be 'info', got %s", config.Logging.Level)
	}

	useTLS, err := strconv.ParseBool(os.Getenv("USE_TLS"))
	if err != nil {
		t.Errorf("Error parsing USE_TLS: %v", err)
	}
	if config.ServerConf.UseTLS != useTLS {
		t.Errorf("Expected UseTLS to be %t, got %t", useTLS, config.ServerConf.UseTLS)
	}

	if config.ServerConf.CertPath != "/app/certs/panel.crt" {
		t.Errorf("Expected CertPath to be '/app/certs/panel.crt', got %s", config.ServerConf.CertPath)
	}

	if config.ServerConf.KeyPath != "/app/certs/panel.key" {
		t.Errorf("Expected KeyPath to be '/app/certs/panel.key', got %s", config.ServerConf.KeyPath)
	}
}
