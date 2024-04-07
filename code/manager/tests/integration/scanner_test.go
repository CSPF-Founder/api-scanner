package integration

import (
	"context"
	"testing"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/manager/db/models"
	"github.com/CSPF-Founder/api-scanner/code/manager/enums/jobstatus"
	"github.com/CSPF-Founder/api-scanner/code/manager/internal/scanner"
	"github.com/CSPF-Founder/api-scanner/code/manager/logger"
)

func TestScanner_GetJobsToScan(t *testing.T) {

	if testDB == nil {
		t.Fatalf("Could not connect to database")
	}

	lgr, err := logger.NewLogger(&logger.Config{
		Level:    "debug",
		Filename: "/tmp/manager.log",
	})
	if err != nil {
		t.Fatalf("Could not create logger: %s", err)
	}

	dbService := models.New(testDB)
	jobService := &dbService.Job
	scanner := scanner.NewScanner("scanner", jobService, lgr)

	ctx := context.Background()
	jobs := scanner.GetJobsToScan(ctx)

	if len(jobs) != 0 {
		t.Fatalf("Expected 0 jobs, got %d", len(jobs))
	}

	createdAt := time.Now()

	// Create a new job
	job := models.Job{
		ID:        1,
		Status:    int(jobstatus.Default),
		APIURL:    "http://localhost:8080",
		CreatedAt: createdAt,
		UserID:    1,
	}

	// Insert the job
	_, err = jobService.Create(ctx, job)

	if err != nil {
		t.Fatalf("Could not insert job: %s", err)
	}

	jobs = scanner.GetJobsToScan(ctx)
	if len(jobs) != 1 {
		t.Fatalf("Expected 1 job, got %d", len(jobs))
	}
}
