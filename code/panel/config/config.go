package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CSPF-Founder/api-scanner/code/panel/logger"
	"github.com/joho/godotenv"
)

// Server represents the Server configuration details
type ServerConfig struct {
	ServerAddress        string   `json:"server_address"`
	CSRFKey              string   `json:"csrf_key"`
	UseTLS               bool     `json:"use_tls"`
	AllowedInternalHosts []string `json:"allowed_internal_hosts"`
	TrustedOrigins       []string `json:"trusted_origins"`
	CertPath             string   `json:"cert_path"`
	KeyPath              string   `json:"key_path"`
	CSRFName             string   `json:"csrf_name"`
}

type DatabaseConfig struct {
	Host   string `json:"host"`
	User   string `json:"user"`
	Pass   string `json:"password"`
	DBName string `json:"db_name"`
}

// Config represents the configuration information.
type Config struct {
	ServerConf             ServerConfig   `json:"server"`
	DBMSType               string         `json:"dbms_type"`
	DatabaseURI            string         `json:"database_uri"`
	DBSSLCaPath            string         `json:"db_sslca_path"`
	MigrationsPath         string         `json:"migrations_prefix"`
	TestFlag               bool           `json:"test_flag"`
	ContactAddress         string         `json:"contact_address"`
	Logging                *logger.Config `json:"logging"`
	ProductTitle           string         `json:"product_title"`
	CopyrightFooterCompany string         `json:"copyright_footer_company"`
	WorkDir                string         `json:"work_dir"`
	TempUploadsDir         string         `json:"temp_uploads_dir"`
}

// Version contains the current project version
var Version = "1"

// ServerName is the server type that is returned in the transparency response.
const ServerName = "api_scanner"

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

func getEnvValueOrError(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.GetFallBackLogger().Error(fmt.Sprintf("Environment variable %s not set", key), nil)
		os.Exit(1)
	}
	return value
}

// LoadConfig loads the configuration from the specified filepath
func LoadConfig() (*Config, error) {
	if os.Getenv("USE_DOTENV") != "false" {
		loadEnv()
	}

	useTLSEnv := os.Getenv("USE_TLS")
	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")
	allowedInternalHostsEnv := os.Getenv("ALLOWED_INTERNAL_HOSTS") // comma separated list of hosts
	trustedOriginsEnv := os.Getenv("TRUSTED_ORIGINS")

	allowedInternalHosts := []string{}
	if allowedInternalHostsEnv != "" {
		for _, host := range strings.Split(allowedInternalHostsEnv, ",") {
			allowedInternalHosts = append(allowedInternalHosts, strings.TrimSpace(host))
		}
	}

	trustedOrigins := []string{}
	if trustedOriginsEnv != "" {
		for _, host := range strings.Split(trustedOriginsEnv, ",") {
			trustedOrigins = append(trustedOrigins, strings.TrimSpace(host))
		}
	}

	srvConfig := ServerConfig{
		ServerAddress:        getEnvValueOrError("SERVER_ADDRESS"),
		CSRFKey:              os.Getenv("CSRF_KEY"),
		CSRFName:             "csrf_token",
		UseTLS:               useTLSEnv == "true", // default to false
		AllowedInternalHosts: allowedInternalHosts,
		TrustedOrigins:       trustedOrigins,
		CertPath:             certPath,
		KeyPath:              keyPath,
	}

	config := &Config{
		ServerConf:     srvConfig,
		DBMSType:       getEnvValueOrError("DBMS_TYPE"),
		DatabaseURI:    getEnvValueOrError("DATABASE_URI"),
		DBSSLCaPath:    os.Getenv("DB_SSLCA_PATH"),
		MigrationsPath: getEnvValueOrError("MIGRATIONS_PREFIX"),
		TestFlag:       false,
		ContactAddress: os.Getenv("CONTACT_ADDRESS"),
		Logging: &logger.Config{
			Level:    os.Getenv("LOG_LEVEL"),
			Filename: getEnvValueOrError("LOG_FILENAME"),
		},
		ProductTitle:           getEnvValueOrError("PRODUCT_TITLE"),
		CopyrightFooterCompany: getEnvValueOrError("COPYRIGHT_FOOTER_COMPANY"),
		WorkDir:                getEnvValueOrError("WORK_DIR"),
		TempUploadsDir:         getEnvValueOrError("TEMP_UPLOADS_DIR"),
	}

	return config, nil
}
