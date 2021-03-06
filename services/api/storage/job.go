package storage

import (
	"fmt"
	"github.com/zale144/instagram-bot/services/api/model"
	"github.com/jinzhu/gorm"
)

type JobStorage struct{}

// GetAll returns all the jobs
func (j JobStorage) GetAll() ([]model.Job, error) {
	var jobs []model.Job
	err := model.DB.Order("id desc").
		Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetByHashTag fetches a job by hashtag
func (j JobStorage) GetByHashTag(hashtag string) (*model.Job, error) {
	var job model.Job
	err := model.DB.
		Where(model.Job{HashTagName: hashtag}).
		First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetOngoingByHashTag fetches an ongoing job by hashtag
func (j JobStorage) GetOngoingByHashTag(hashtag string) (*model.Job, error) {
	var job model.Job
	err := model.DB.
		Where(model.Job{HashTagName: hashtag}).Where("finished_at=?", 0).
		First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// Insert saves a job to the database
func (j JobStorage) Insert(job *model.Job) error {
	tx := model.DB.Begin()
	if err := tx.Save(job).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Update a job
type JobUpdater struct {
	jobID   uint
	updates map[string]interface{}
}

// NewJobUpdater creates an instance of a JobUpdater struct
func (p JobStorage) NewJobUpdater(jobID uint) *JobUpdater {
	return &JobUpdater{
		jobID:   jobID,
		updates: make(map[string]interface{}),
	}
}

// FinishedAt sets the 'finished_at' column for updating
func (a *JobUpdater) FinishedAt(f int64) *JobUpdater {
	a.updates["finished_at"] = f
	return a
}

// Update commits the update request
func (a *JobUpdater) Update(tx *gorm.DB) error {
	if tx == nil {
		tx = model.DB
	}
	tx = tx.Model(&model.Job{Model: gorm.Model{ID: a.jobID}}).
		Updates(a.updates)
	rowsAffected, err := tx.RowsAffected, tx.Error

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		err = fmt.Errorf("record not found")
		return err
	}
	return nil
}
