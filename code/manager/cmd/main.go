package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/manager/config"
	"github.com/CSPF-Founder/api-scanner/code/manager/db"
	"github.com/CSPF-Founder/api-scanner/code/manager/db/models"
	"github.com/CSPF-Founder/api-scanner/code/manager/internal/scanner"
	"github.com/CSPF-Founder/api-scanner/code/manager/logger"
	"github.com/CSPF-Founder/api-scanner/code/manager/utils"
)

type application struct {
	Config config.AppConfig
	DB     models.Service
	logger *logger.Logger
}

func main() {

	conf := config.LoadConfig()
	appLogger, err := logger.NewLogger(conf.Logging)
	if err != nil {
		log.Fatal("Error setting up logging: ", err)
	}
	appLogger.Info("Initializing API Scanner Manager...")

	conn, err := db.ConnectDBWithRetry(conf.DatabaseURI, 3, 5*time.Second)
	if err != nil {
		appLogger.Fatal("Unable to connect to database", err)
	}
	appLogger.Info("Connected to database")
	defer conn.Close()

	// Wrapper for the SQLC generated models
	app := &application{
		Config: conf,
		DB:     models.New(conn),
		logger: appLogger,
	}

	scannerInstance := scanner.NewScanner(conf.ScannerCmd, &app.DB.Job, app.logger)

	appLogger.Info("Service is Running...")

	ctx, cancel := context.WithCancel(context.Background())
	// Set up a signal channel to capture interrupt and termination signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		// Wait for the interrupt signal
		<-interrupt

		// Perform cleanup operations before exiting (if needed)
		appLogger.Info("Service is stopping...")

		// Cancel the context to signal a graceful shutdown
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			appLogger.Info("Shutting down gracefully...")
			return
		default:

			scannerInstance.ProcessScanQueue(ctx)

			if err := utils.SleepContext(ctx, 20*time.Second); err != nil {
				appLogger.Error("Error sleeping", err)
			}
		}
	}
}
