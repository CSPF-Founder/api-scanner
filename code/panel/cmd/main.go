package main

import (
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/CSPF-Founder/api-scanner/code/panel/config"
	"github.com/CSPF-Founder/api-scanner/code/panel/controllers"
	"github.com/CSPF-Founder/api-scanner/code/panel/logger"
)

func main() {
	// Load the config
	conf, err := config.LoadConfig()
	// Just warn if a contact address hasn't been configured
	if err != nil {
		log.Fatal("Error loading config", err)
	}

	appLogger, err := logger.NewLogger(conf.Logging)
	if err != nil {
		log.Fatal("Error setting up logging", err)
	}

	app, err := baseSetup(conf, appLogger)
	if err != nil {
		appLogger.Fatal("Error setting up app", err)
		return
	}

	//Start Server
	go app.StartServer()

	// Handle shutdown here
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	appLogger.Info("CTRL+C Received... shutting down servers")
	defer func() {
		_ = app.Shutdown()
	}()
}

func baseSetup(conf *config.Config, appLogger *logger.Logger) (*controllers.App, error) {
	// create directory if not exists
	if _, err := os.Stat(conf.TempUploadsDir); errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(conf.TempUploadsDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	app := controllers.NewApp(conf, appLogger)
	app.SetupDB()
	app.SetupRoutes()

	return app, nil
}
