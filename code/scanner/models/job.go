package models

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/enums/jobstatus"
)

const TABLE_NAME string = "jobs"

type Job struct {
	ID            int
	UserID        int
	Status        jobstatus.JobStatus
	ScannerID     int
	ScannerIP     string
	CreatedAt     time.Time
	CompletedTime *time.Time
}

// Model Functions
func (jm *Job) GetLocalWorkDir(cfg config.Config) string {
	return filepath.Join(
		cfg.LocalTempDir,
		fmt.Sprintf("user_%d", jm.UserID),
		fmt.Sprintf("job_%d", jm.ID),
	)
}

func (jm *Job) GetRemoteWorkDir(cfg config.Config) string {
	return filepath.Join(
		cfg.RemoteWorkDir,
		fmt.Sprintf("user_%d", jm.UserID),
		fmt.Sprintf("job_%d", jm.ID),
	)
}

// CRUD operations for Job model:

func (m *DBModel) UpdateJobStatus(ctx context.Context, jobID int, jobStatus jobstatus.JobStatus) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	query := "UPDATE " + TABLE_NAME + " SET status=? WHERE id=?"

	_, err := m.DB.QueryContext(ctx, query, int(jobStatus), int(jobID))
	if err == nil {
		return true, nil
	}
	return false, err
}

func (m *DBModel) MarkAsCompleted(ctx context.Context, jobID int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := "UPDATE " + TABLE_NAME + " SET status=?, completed_time=? WHERE id=?"
	_, err := m.DB.QueryContext(ctx, query, jobstatus.ScanCompleted, time.Now(), jobID)
	if err == nil {
		return true, nil
	}
	return false, err
}

func (m *DBModel) Add(ctx context.Context, job *Job) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := "INSERT INTO " + TABLE_NAME + " (user_id, status) VALUES (?, ?)"
	_, err := m.DB.QueryContext(ctx, query, job.UserID, job.Status)
	if err == nil {
		return true, nil
	}
	return false, err
}
