package models

import (
	"regexp"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/panel/enums/jobstatus"
	"github.com/DATA-DOG/go-sqlmock"
)

func (ctx *testContext) mockJobByID(job Job) {
	jobRow := sqlmock.NewRows([]string{"id", "status", "created_at", "user_id"}).
		AddRow(job.ID, job.Status, job.CreatedAt, job.UserID)

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE id=? ORDER BY `jobs`.`id` LIMIT 1")).
		WithArgs(job.ID).
		WillReturnRows(jobRow)
}

// func (ctx *testContext) mockJobByIDWithUser(job Job) {
// 	jobRow := sqlmock.NewRows([]string{"id", "status", "created_at", "user_id"}).
// 		AddRow(job.ID, job.Status, job.CreatedAt, job.UserID)

//		ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE id=? AND user_id=? ORDER BY `jobs`.`id` LIMIT 1")).
//			WithArgs(job.ID, job.UserID).
//			WillReturnRows(jobRow)
//	}
func TestGetJob(t *testing.T) {
	ctx := setupTest(t)
	testJob := Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}
	ctx.mockJobByID(testJob)
	got, err := GetJobByID(testJob.ID)
	if err != nil {
		t.Fatalf("error getting job: %v", err)
	}

	if got.ID != testJob.ID {
		t.Fatalf("invalid job ID. expected %d got %d", testJob.ID, got.ID)
	}

	if got.Status != testJob.Status {
		t.Fatalf("invalid job status. expected %d got %d", testJob.Status, got.Status)
	}

}

func TestGetJobWithUser(t *testing.T) {
	ctx := setupTest(t)
	testJob := Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}
	ctx.mockJobByID(testJob)
	got, err := GetByIDAndUser(testJob.ID, testJob.UserID)
	if err != nil {
		t.Fatalf("error getting job: %v", err)
	}

	if got.ID != testJob.ID {
		t.Fatalf("invalid job ID. expected %d got %d", testJob.ID, got.ID)
	}

	if got.Status != testJob.Status {
		t.Fatalf("invalid job status. expected %d got %d", testJob.Status, got.Status)
	}

}
func TestDeleteJob(t *testing.T) {
	ctx := setupTest(t)
	testJob := Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}

	// Mock job deletion
	ctx.mock.ExpectBegin()
	ctx.mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `jobs` WHERE `jobs`.`id` = ?")).
		WithArgs(testJob.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ctx.mock.ExpectCommit()

	err := DeleteJob(testJob.ID)
	if err != nil {
		t.Fatalf("error deleting job: %v", err)
	}

	if err := ctx.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetJobs(t *testing.T) {
	ctx := setupTest(t)
	testUser := User{ID: 1}
	testJobs := []Job{
		{ID: 1, Status: jobstatus.ScanStarted, UserID: testUser.ID},
		{ID: 2, Status: jobstatus.ScanCompleted, UserID: testUser.ID},
	}

	jobsRow := sqlmock.NewRows([]string{"id", "status", "user_id"})
	for _, job := range testJobs {
		jobsRow = jobsRow.AddRow(job.ID, job.Status, job.UserID)
	}

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE user_id = ?")).
		WithArgs(testUser.ID).
		WillReturnRows(jobsRow)

	got, err := GetJobs(&testUser)
	if err != nil {
		t.Fatalf("error getting jobs: %v", err)
	}

	if len(got) != len(testJobs) {
		t.Fatalf("expected %d jobs, got %d", len(testJobs), len(got))
	}

	for i, job := range got {
		if job.ID != testJobs[i].ID || job.Status != testJobs[i].Status {
			t.Errorf("mismatch in job at index %d, expected %+v, got %+v", i, testJobs[i], job)
		}
	}
}

// TestSaveJob tests saving a job
func TestSaveJob(t *testing.T) {
	ctx := setupTest(t)
	testJob := Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}

	ctx.mock.ExpectBegin()
	ctx.mock.ExpectExec(regexp.QuoteMeta("UPDATE `jobs` SET `status`=?,`api_url`=?,`created_at`=?,`completed_time`=?,`user_id`=? WHERE `id` = ?")).
		WithArgs(testJob.Status, testJob.ApiURL, testJob.CreatedAt, testJob.CompletedTime, testJob.UserID, testJob.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	ctx.mock.ExpectCommit()

	err := SaveJob(&testJob)
	if err != nil {
		t.Fatalf("error saving job: %v", err)
	}

	if err := ctx.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
