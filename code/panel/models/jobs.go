package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/CSPF-Founder/api-scanner/code/panel/enums/jobstatus"
)

const DateTimeFormat = "2006-01-02 3:4 PM"

// Job represents the job model.
type Job struct {
	ID                  uint64              `gorm:"default:uuid_short()"`
	Status              jobstatus.JobStatus `json:"status" sql:"not null"`
	StatusText          string              `gorm:"-"`
	ApiURL              string              `json:"api_url" sql:"null"`
	CreatedAt           time.Time           `json:"created_at" sql:"null"`
	CompletedTime       time.Time           `json:"completed_time" sql:"null"`
	CompletedTimeString string              `gorm:"-"`
	UserID              int64               `json:"user_id" sql:"not null"`
	ScanCompleted       bool                `gorm:"-"`
}

func (j *Job) AfterFind(tx *gorm.DB) (err error) {
	j.StatusText = j.Status.GetText()
	j.CompletedTimeString = j.CompletedTime.Format(DateTimeFormat)
	j.ScanCompleted = j.Status == jobstatus.ScanCompleted
	return
}

// GetJobByID returns the job that the given id corresponds to.
// If no job is found, an error is returned.
func GetJobByID(id uint64) (Job, error) {
	var job Job
	err := db.Where("id=?", id).First(&job).Error
	return job, err
}

// GetByIDAndUser returns the job that the given id corresponds to.
// If no job is found, an error is returned.
func GetByIDAndUser(id uint64, userID int64) (Job, error) {
	var job Job
	err := db.Where("id=? AND user_id=?", id, userID).First(&job).Error
	return job, err
}

// Delete Job removes the job from table
// error is thrown if anything happens.
func DeleteJob(id uint64) error {
	err := db.Delete(&Job{}, id).Error
	return err
}

// GetJobs returns the jobs
func GetJobs(u *User) ([]Job, error) {
	us := []Job{}
	err := db.Where("user_id = ?", u.ID).Find(&us).Error
	return us, err
}

// SaveJob saves the job to the database
func SaveJob(u *Job) error {
	return db.Save(&u).Error
}
