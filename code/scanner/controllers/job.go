package controllers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/enums/jobstatus"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/filesync"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/scanner"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
	"github.com/CSPF-Founder/api-scanner/code/scanner/logger"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
	"github.com/CSPF-Founder/api-scanner/code/scanner/pkg/openapi"
	"github.com/go-ini/ini"
)

// JobController is a struct that contains the configuration and the job model
type JobController struct {
	config      config.Config
	logger      *logger.FileLogger
	job         models.Job
	scanDir     string
	logFilePath string
}

// constructor for JobController
func NewJobController(config config.Config, jobID int, userID int) (*JobController, error) {

	job := models.Job{
		ID:     jobID,
		UserID: userID,
	}

	scanDir := job.GetLocalWorkDir(config)

	if _, err := os.Stat(scanDir); os.IsNotExist(err) {
		err := os.MkdirAll(scanDir, os.ModePerm) // os.ModePerm is 777, might need to check
		if err != nil {
			return nil, fmt.Errorf("Error creating scan directory")
		}
	}

	return &JobController{
		config:      config,
		job:         job,
		scanDir:     scanDir,
		logFilePath: filepath.Join(scanDir, "logs"),
	}, nil
}

// UpdateJobStatus updates the status of the job
func (c *JobController) UpdateJobStatus(ctx context.Context, model models.DBModel, jobStatus jobstatus.JobStatus) {
	isUpdated, err := model.UpdateJobStatus(ctx, c.job.ID, jobStatus)
	if isUpdated {
		c.logger.Info("Job Status Updated")
	} else {
		c.logger.Error("Job Status Failed to Update", err)

	}
}

// MarkAsFinished marks the job as finished
func (c *JobController) MarkAsFinished(ctx context.Context, model models.DBModel, jobID int) {
	isUpdated, err := model.MarkAsCompleted(ctx, c.job.ID)
	if isUpdated {
		c.logger.Info("Job Marked Finished")
	} else {
		c.logger.Error("Job Marked Failed to Finish", err)

	}
}

// MarkAsFailed marks the job as failed
func (c *JobController) handleScanFailure(
	ctx context.Context,
	dbModel models.DBModel,
	err error,
	errorStatus jobstatus.JobStatus,
) error {
	if c.logger != nil {
		c.logger.Error("Error running scan", err)
	}

	if _, err := os.Stat(c.logFilePath); err == nil {
		err = filesync.CopyToRemoteOnlyLog(c.config, &c.job)
		if err != nil {
			return fmt.Errorf("Error copying files to remote %v: %v", c.job.ID, err)
		}
	}

	c.UpdateJobStatus(ctx, dbModel, errorStatus)
	return nil
}

// getAuthHeaders returns the headers from the auth header file
func (c *JobController) getAuthHeaders(authHeaderFile string) ([]schemas.AuthHeaderMap, error) {
	// Check if auth file has headers
	authConfig, err := ini.Load(authHeaderFile)
	if err != nil {
		c.logger.Fatal("Error loading auth header file", err)
		return nil, err
	}

	loadedHeaders := authConfig.Section("AUTH_HEADERS").Keys()
	numberofHeaders := len(loadedHeaders)
	if numberofHeaders == 0 {
		c.logger.Info("No Security Headers")
		return nil, nil
	} else if numberofHeaders > 1 {
		c.logger.Error("Multiple Security Headers", nil)
		return nil, nil
	}

	headerData := []schemas.AuthHeaderMap{}
	if numberofHeaders == 1 {
		// return []string{authHeaders[0].String()}, nil
		firstHeader := schemas.AuthHeaderMap{
			Name:  loadedHeaders[0].Name(),
			Value: loadedHeaders[0].String(),
		}
		headerData = append(headerData, firstHeader)
		return headerData, nil
	}

	return nil, fmt.Errorf("Invalid Auth Header file")
}

// Run function handles the job and performs the scan
func (c *JobController) Run(ctx context.Context, model models.DBModel) error {

	authHeaderFile := filepath.Join(c.scanDir, "auth_headers.conf")
	openApiFile := filepath.Join(c.scanDir, "openapi.yaml")

	_, err := filesync.CopyFromRemote(c.config, c.logger, &c.job)
	if err != nil {
		return c.handleScanFailure(ctx, model, fmt.Errorf("Input files not present in expected path"), jobstatus.ScanFailed)
	}

	// Scan Log Setup:
	scanLogFile, err := os.OpenFile(c.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer scanLogFile.Close()

	// Create a logger for the scanner
	c.logger, err = logger.GetFileLogger(scanLogFile)
	if err != nil {
		return err
	} else if c.logger == nil {
		return fmt.Errorf("Error creating scanner logger")
	}

	c.logger.Info("Scan Initiated")

	_, err = os.Stat(openApiFile)
	if err != nil && os.IsNotExist(err) {
		return c.handleScanFailure(
			ctx, model,
			fmt.Errorf("OpenAPI file not present in expected path"),
			jobstatus.ScanFailed,
		)
	}
	_, authHeaderFileExistError := os.Stat(authHeaderFile)
	if authHeaderFileExistError != nil && os.IsNotExist(authHeaderFileExistError) {
		return c.handleScanFailure(
			ctx,
			model,
			fmt.Errorf("authHeaderFile file not present in expected path"), jobstatus.ScanFailed,
		)
	}

	authHeaderData, err := c.getAuthHeaders(authHeaderFile)
	if err != nil {
		return err
	}

	serverURL, err := openapi.GetServerURLFromOpenAPI(openApiFile)
	if err != nil {
		return c.handleScanFailure(
			ctx,
			model,
			fmt.Errorf("Invalid OpenAPI: %v", err),
			jobstatus.ScanFailed,
		)
	}

	scanner := scanner.Scanner(c.job, c.logger)
	scanCompleted := scanner.Run(
		ctx,
		c.config,
		model,
		serverURL,
		openApiFile,
		authHeaderData,
	)

	if scanCompleted {
		err := c.RunReporter(ctx, &model)
		if err != nil {
			return c.handleScanFailure(
				ctx,
				model,
				fmt.Errorf("Reporter failed %v", err),
				jobstatus.ScanFailed,
			)
		}
		return nil
	}

	c.logger.Info("Unable to finish scan")

	return c.handleScanFailure(
		ctx,
		model,
		errors.New("Unable to finish scan"),
		jobstatus.ScanFailed,
	)

}

func (c *JobController) RunReporter(ctx context.Context, model *models.DBModel) error {
	scriptArgs := []string{"-m", "reporter", "-j", strconv.Itoa(c.job.ID), "-u", strconv.Itoa(c.job.UserID)}
	cmd := exec.Command(c.config.ReporterBinPath, scriptArgs...)

	// Redirect standard output and error to the Go program's standard output and error
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return err
	}

	_, copyerr := filesync.CopyToRemote(c.config, &c.job)
	if copyerr != nil {
		c.logger.Error("Error copying files to remote", copyerr)
		return err
	}

	c.UpdateJobStatus(ctx, *model, jobstatus.FilesCopiedToRemote)
	c.MarkAsFinished(ctx, *model, c.job.ID)

	return nil
}
