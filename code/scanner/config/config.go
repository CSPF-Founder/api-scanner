package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/scanner/logger"
	"github.com/joho/godotenv"
)

// Load loads the environment variables from the .env file
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// AppConfig is the struct for the application configuration
type Config struct {
	Env             string
	DSN             string
	RemoteWorkDir   string
	LogLevel        string
	LocalTempDir    string
	ScannerImage    string
	ReporterBinPath string
	Logging         *logger.Config
	ScanTimeout     time.Duration
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

// func getEnvValueOrError(key string) string {
// 	value := os.Getenv(key)
// 	if value == "" {
// 		logger.Error(fmt.Sprintf("Environment variable %s not set", key), nil)
// 		os.Exit(1)
// 	}
// 	return value
// }

// Load loads the environment variables from the .env file
func Load() *Config {
	if os.Getenv("USE_DOTENV") != "false" {
		loadEnv()
	}

	dsn := getEnv("DSN", "")
	if dsn == "" {
		log.Fatal("DSN", "")
		log.Fatal("DSN is not specified in the environment")
	}

	config := &Config{
		Env:             getEnv("ENVIRONMENT", "dev"),
		RemoteWorkDir:   getEnv("REMOTE_WORK_DIR", ""),
		LocalTempDir:    getEnv("LOCAL_TEMP_DIR", ""),
		ScannerImage:    getEnv("SCANNER_IMAGE", ""),
		ReporterBinPath: getEnv("REPORTER_BIN_PATH", ""),
		DSN:             dsn,
		Logging: &logger.Config{
			Level: os.Getenv("LOG_LEVEL"),
		},
		ScanTimeout: 2 * time.Hour,
	}

	return config
}
