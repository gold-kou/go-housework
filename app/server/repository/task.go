package repository

import (
	"errors"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/jinzhu/gorm"
)

// TaskRepositoryInterface is a repository interface of task
type TaskRepositoryInterface interface {
	InsertTask(*db.Task) error
	SelectTaskWhereFamilyIDDate(uint64, string) ([]*db.Task, error)
	UpdateTask(*db.Task) error
	DeleteTaskWhereMemberID(*db.Task) error
}

// TaskRepository is a repository of task
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a pointer of TaskRepository
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// InsertTask insert task
func (r TaskRepository) InsertTask(task *db.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// SelectTaskWhereFamilyIDDate select task where family_id & date
func (r TaskRepository) SelectTaskWhereFamilyIDDate(familyID uint64, targetDate string) ([]*db.Task, error) {
	var dbTasks []*db.Task
	if err := r.db.Where("family_id = ? AND date = ?", familyID, targetDate).Find(&dbTasks).Error; err != nil {
		return dbTasks, common.NewInternalServerError(err.Error())
	}
	return dbTasks, nil
}

// UpdateTask update task
func (r TaskRepository) UpdateTask(task *db.Task) error {
	if err := r.db.Save(task).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// DeleteTaskWhereMemberID delete task
func (r TaskRepository) DeleteTaskWhereMemberID(task *db.Task) error {
	result := r.db.Where("member_id = ?", task.MemberID).Delete(task)
	if result.RowsAffected == 0 {
		return errors.New("delete failed")
	}
	return nil
}
