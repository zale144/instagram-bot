package storage

import (
	"fmt"

	"github.com/zale144/instagram-bot/api/model"

	"github.com/jinzhu/gorm"
)

type ProcessedUserStorage struct{}

// GetAll returns all the users
func (pus ProcessedUserStorage) GetAll(page uint) ([]model.ProcessedUser, error) {
	// page size is 10
	offset := (page - 1) * 10
	limit := 10
	var users []model.ProcessedUser
	err := model.DB.Order("id desc").Preload("Job").
		Offset(offset). // skip (page * 10) number of rows
		Limit(limit).   // limit to 10 rows
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByJob returns all the users by job id
func (pus ProcessedUserStorage) GetByJob(jobID, page uint) ([]model.ProcessedUser, error) {
	// page size is 10
	offset := (page - 1) * 10
	limit := 10
	var users []model.ProcessedUser
	err := model.DB.Order("id desc").Preload("Job").
		Where("job_id=?", jobID).
		Offset(offset). // skip (page * 10) number of rows
		Limit(limit).   // limit to 10 rows
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetOne fetches a user by id
func (pus ProcessedUserStorage) GetOne(id uint) (*model.ProcessedUser, error) {
	var user model.ProcessedUser
	err := model.DB.Order("id desc").Preload("Job").Where("id=?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername fetches a user by username
func (pus ProcessedUserStorage) GetByUsername(username string) (*model.ProcessedUser, error) {
	var user model.ProcessedUser
	err := model.DB.Order("id desc").Where("username=?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Insert saves a processed user to the database
func (pus ProcessedUserStorage) Insert(user model.ProcessedUser) error {
	tx := model.DB.Begin()
	if err := tx.Save(&user).Omit("job").Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// Update a processed user
type ProcessedUserUpdater struct {
	ProcessedUserID uint
	updates         map[string]interface{}
}

// NewProcessedUserUpdater returns a new ProcessedUserUpdater struct
func (p ProcessedUserStorage) NewProcessedUserUpdater(userID uint) *ProcessedUserUpdater {
	return &ProcessedUserUpdater{
		ProcessedUserID: userID,
		updates:         make(map[string]interface{}),
	}
}

// Successful sets the updater for the 'successful' column
func (a *ProcessedUserUpdater) Successful(f int64) *ProcessedUserUpdater {
	a.updates["successful"] = f
	return a
}

// Update commits the updates to the database
func (a *ProcessedUserUpdater) Update(tx *gorm.DB) error {
	if tx == nil {
		tx = model.DB
	}
	tx = tx.Model(&model.ProcessedUser{Model: gorm.Model{ID: a.ProcessedUserID}}).
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
