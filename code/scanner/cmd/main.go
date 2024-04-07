package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/controllers"
	"github.com/CSPF-Founder/api-scanner/code/scanner/db"
	"github.com/CSPF-Founder/api-scanner/code/scanner/logger"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
)

type application struct {
	Config *config.Config
	DB     models.DBModel
}

type CLIInput struct {
	Module string
	JobID  int
	UserID int
}

func main() {
	logger.GetFallBackLogger().Info("Starting API Scanner...")

	conf := config.Load()
	appLogger, err := logger.NewLogger(conf.Logging)
	if err != nil {
		logger.GetFallBackLogger().Error("Error initializing logger", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), conf.ScanTimeout)
	// Set up a signal channel to capture interrupt and termination signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Handle signals in a goroutine
	go func() {
		select {
		case <-interrupt:
			// Wait for the interrupt signal
			appLogger.Info("Stopping scanner...")
		case <-ctx.Done():
			// Context has timed out
			appLogger.Info("Scanner timed out...")
		}
		// Cancel the context to signal a graceful shutdown
		cancel()
	}()

	conn, err := db.ConnectDB(conf.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app := &application{
		Config: conf,
		DB:     models.DBModel{DB: conn},
	}

	cli, err := handleCLI()
	if err != nil {
		appLogger.Error("Error parsing CLI input", err)
		return
	}

	switch cli.Module {
	case "scanner":
		mainController, err := controllers.NewJobController(*app.Config, cli.JobID, cli.UserID)
		if err != nil {
			appLogger.Error("Error running scanner", err)
		}

		err = mainController.Run(ctx, app.DB)
		if err != nil {
			appLogger.Error("Error running scanner", err)
		}
	default:
		flag.PrintDefaults()
	}

}

// handleCLI parses the command line input and returns a CLIInput struct
func handleCLI() (*CLIInput, error) {
	cli := &CLIInput{}

	flag.StringVar(&cli.Module, "m", "", "Module to run (choices: scanner)")
	flag.IntVar(&cli.JobID, "j", 0, "Job ID")
	flag.IntVar(&cli.UserID, "u", 0, "User ID")

	flag.Parse()

	if cli.Module == "" || cli.JobID == 0 || cli.UserID == 0 {
		return nil, errors.New("Invalid input")
	}

	return cli, nil
}
