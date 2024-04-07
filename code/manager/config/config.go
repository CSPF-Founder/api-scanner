package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CSPF-Founder/api-scanner/code/manager/logger"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env           string
	DatabaseURI   string
	ScannerDocker string
	LogLevel      string
	ScannerCmd    string
	Logging       *logger.Config
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getEnvValueOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.GetFallBackLogger().Error(fmt.Sprintf("Environment variable %s not set", key), nil)
		os.Exit(1)
	}
	return value
}

func loadEnv() {
	//determin bin directory and load .env from there
	exe, err := os.Executable()
	if err != nil {
		logger.GetFallBackLogger().Fatal("Error loading .env file", err)
	}
	binDir := filepath.Dir(exe)
	envPath := filepath.Join(binDir, ".env")
	if err := godotenv.Load(envPath); err == nil {
		return
	}

	// try to load .env from current directory
	envPath = ".env"
	if err := godotenv.Load(envPath); err == nil {
		return
	}
	logger.GetFallBackLogger().Error("Error loading .env file", err)
	os.Exit(1)

}

// LoadConfig loads the configuration from the specified filepath
func LoadConfig() AppConfig {
	if os.Getenv("USE_DOTENV") != "false" {
		loadEnv()
	}

	logLevel := getEnv("LOG_LEVEL", "debug")
	logFileName := getEnv("LOG_FILENAME", "logfile")

	return AppConfig{
		DatabaseURI:   getEnvValueOrPanic("DATABASE_URI"),
		ScannerDocker: getEnvValueOrPanic("SCANNER_DOCKER"),
		LogLevel:      logLevel,
		ScannerCmd:    getEnvValueOrPanic("SCANNER_CMD"),
		Logging: &logger.Config{
			Level:    logLevel,
			Filename: logFileName,
		},
	}
}
