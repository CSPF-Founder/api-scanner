package scanner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/manager/db/models"
	"github.com/CSPF-Founder/api-scanner/code/manager/enums/jobstatus"
	"github.com/CSPF-Founder/api-scanner/code/manager/logger"
	"github.com/CSPF-Founder/api-scanner/code/manager/utils"
)

type Scanner struct {
	ScannerCmd string
	JobService *models.JobService
	logger     *logger.Logger
}

func NewScanner(scannerCmd string, jobService *models.JobService, lgr *logger.Logger) *Scanner {
	return &Scanner{
		ScannerCmd: scannerCmd,
		JobService: jobService,
		logger:     lgr,
	}
}

const ScannerTimeout = 2 * time.Hour // 2-hour timeout

// ScanJobs scans the jobs table with status=DEFAULT
// and returns and array of jobs
func (s *Scanner) GetJobsToScan(ctx context.Context) []models.Job {

	jobs, err := s.JobService.GetByStatus(ctx, jobstatus.Default)
	if err != nil {
		s.logger.Error("Failed to find jobs", err)
	}
	return jobs
}

func (s *Scanner) runSubprocess(ctx context.Context, jobID int, userID int) error {
	jobIDStr := fmt.Sprintf("%d", jobID)
	userIDStr := fmt.Sprintf("%d", userID)

	ctx, cancel := context.WithTimeout(ctx, ScannerTimeout)
	defer cancel()

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, s.ScannerCmd, "-j", jobIDStr, "-u", userIDStr, "-m", "scanner")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		s.logger.Info(fmt.Sprintf("Command failed with error: %s", stderr.String()))
		return err
	}

	// s.logger.Info(fmt.Sprintf("Scanner cmd output: %s", out.String()))

	return nil
}

// Updates the job table with status=JobStatusScanFailed
func (s *Scanner) markJobFailed(ctx context.Context, job models.Job) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	job.Status = int(jobstatus.ScanFailed)
	isUpdated, err := s.JobService.UpdateStatus(ctx, job)
	if err != nil {
		return err
	}

	if !isUpdated {
		return fmt.Errorf("Failed to update job status")
	}

	return nil
}

func (s *Scanner) ProcessScanQueue(ctx context.Context) {
	jobs := s.GetJobsToScan(ctx)
	if len(jobs) == 0 {
		if err := utils.SleepContext(ctx, 10*time.Second); err != nil {
			s.logger.Error("ProcessScanQueue - Error sleeping when no jobs", err)
		}
		return
	}

	for _, job := range jobs {
		jobID := job.ID
		jobFinished := false

		// Attempt to run the job using subprocess
		err := s.runSubprocess(ctx, jobID, job.UserID)
		if err == nil {
			jobFinished = true
		} else {
			s.logger.Error(fmt.Sprintf("Scanner for Job ID %d Failed", jobID), nil)
		}

		if !jobFinished {
			if err := s.markJobFailed(ctx, job); err != nil {
				s.logger.Error("Failed to mark job as failed", err)
			}
		}

		if err := utils.SleepContext(ctx, 30*time.Second); err != nil {
			s.logger.Error("ProcessScanQueue - Error sleeping", err)
		}
	}
}
