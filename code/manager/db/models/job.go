package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/manager/enums/jobstatus"
)

const TABLE_NAME string = "jobs"

type Job struct {
	ID            int
	Status        int
	UserID        int
	APIURL        string
	CreatedAt     time.Time
	CompletedTime *time.Time
}

type JobService struct {
	DB *sql.DB
}

func NewJobService(db *sql.DB) JobService {
	return JobService{DB: db}
}

const getJobsByStatusSQL = `SELECT
    id,
    status,
    api_url,
    created_at,
    completed_time,
    user_id
FROM jobs
WHERE status = ? ORDER BY created_at ASC
`

// Get jobs by status and returns
// Array of jobs and error
func (s JobService) GetByStatus(ctx context.Context, status jobstatus.JobStatus) ([]Job, error) {
	if s.DB == nil {
		return nil, fmt.Errorf("DB is nil")
	}
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	rows, err := s.DB.QueryContext(ctx, getJobsByStatusSQL, status)
	if err != nil {
		return nil, err
	}
	var jobs []Job

	var createdAtString string
	var completedTime sql.NullTime
	for rows.Next() {
		var job Job
		if err := rows.Scan(
			&job.ID,
			&job.Status,
			&job.APIURL,
			&createdAtString,
			&completedTime,
			&job.UserID,
		); err != nil {
			return nil, err
		}

		if createdAtString != "" {

			createdAt, err := ParseTimeString(createdAtString)
			if err != nil {
				return nil, err
			}
			job.CreatedAt = createdAt
		}

		if completedTime.Valid {
			job.CompletedTime = &completedTime.Time
		}

		jobs = append(jobs, job)
	}
	return jobs, err
}

const updateJobStatusSQL = `UPDATE jobs SET status=? WHERE id=?`

// Updates job status by ID and returns a
// boolean and error
func (s JobService) UpdateStatus(ctx context.Context, job Job) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// query := fmt.Sprintf("UPDATE %s SET status=? WHERE id=?", TABLE_NAME)  // can be used as well not sure about SQL injection
	_, err := s.DB.QueryContext(ctx, updateJobStatusSQL, job.Status, job.ID)
	if err != nil {
		return false, err
	}

	return true, err
}

const createJobSQL = `INSERT INTO jobs (status, user_id, api_url, created_at) VALUES (?, ?, ?, ?)`

// create job
func (s JobService) Create(ctx context.Context, job Job) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, createJobSQL, job.Status, job.UserID, job.APIURL, job.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}
