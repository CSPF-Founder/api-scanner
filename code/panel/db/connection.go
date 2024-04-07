package db

import (
	"crypto/x509"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/CSPF-Founder/api-scanner/code/panel/config"
	"github.com/CSPF-Founder/api-scanner/code/panel/internal/sessions"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var conf *config.Config

const MaxDatabaseConnectionAttempts int = 10

// Copy of auth.GenerateSecureKey to prevent cyclic import with auth library
// func generateSecureKey() string {
// 	k := make([]byte, 32)
// 	_, err := io.ReadFull(rand.Reader, k)
// 	if err != nil {
// 		return ""
// 	}
// 	return fmt.Sprintf("%x", k)
// }

// func chooseDBDriver(name, openStr string) goose.DBDriver {
// 	d := goose.DBDriver{Name: name, OpenStr: openStr}

// 	d.Import = "github.com/go-sql-driver/mysql"
// 	d.Dialect = &goose.MySqlDialect{}

// 	return d
// }

// Setup initializes the database and runs any needed migrations.
//
// First, it establishes a connection to the database, then runs any migrations
// newer than the version the database is on.
//
// Once the database is up-to-date, we create an admin user (if needed) that
// has a randomly generated API key and password.
func Setup(c *config.Config, session *sessions.SessionManager) (err error) {
	// Setup the package-scoped config
	conf = c
	// Setup the goose configuration
	// migrateConf := &goose.DBConf{
	// 	MigrationsDir: conf.MigrationsPath,
	// 	Env:           "production",
	// 	Driver:        chooseDBDriver(conf.DBMSType, conf.DatabaseURI),
	// }
	// // Get the latest possible migration
	// latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	// if err != nil {
	// 	return fmt.Errorf("Error getting latest migration: %s", err)
	// }

	var db *gorm.DB

	switch conf.DBMSType {
	case "mysql":
		if conf.DBSSLCaPath != "" {
			// Register certificates for tls encrypted db connections
			rootCertPool := x509.NewCertPool()
			pem, err := os.ReadFile(conf.DBSSLCaPath)
			if err != nil {
				return fmt.Errorf("Failed to read CA certificate: %s", err)
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				return fmt.Errorf("Failed to append PEM.")
			}

			// err = mysql.RegisterTLSConfig("ssl_ca", &tls.Config{
			// 	RootCAs: rootCertPool,
			// })
			// if err != nil {
			// 	return fmt.Errorf("Failed to register TLS config: %s", err)
			// }
		}

		// Open our database connection
		i := 0
		for {
			db, err = gorm.Open(mysql.Open(conf.DatabaseURI), &gorm.Config{})
			if err == nil {
				break
			}

			if i >= MaxDatabaseConnectionAttempts {
				return fmt.Errorf("Error opening database: %s", err)
			}
			i += 1

			// appLogger.Warn("waiting for database to be up...", nil)
			time.Sleep(5 * time.Second)
		}
	default:
	}

	if db == nil {
		return fmt.Errorf("Unable to open database connection")
	}

	if session.Store, err = gormstore.New(db); err != nil {
		log.Fatal(err)
	}

	// db.LogMode(false)
	// db.SetLogger(log.Logger)
	// db.DB().SetMaxOpenConns(1)

	// Migrate up to the latest version
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect(c.DBMSType); err != nil {
		return fmt.Errorf("Error setting dialect: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		// Handle error
		return fmt.Errorf("Error obtaining *sql.DB from *gorm.DB: %s", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("Error running migrations: %s", err)
	}

	// err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, db.DB())
	// if err != nil {
	// 	return fmt.Errorf("Error running migrations: %s", err)
	// }

	models.SetupDB(db)

	return nil
}
