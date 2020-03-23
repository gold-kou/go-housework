package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// DeleteTaskServiceInterface is a service interface of deleteTask
type DeleteTaskServiceInterface interface {
	Execute() error
}

// DeleteTask struct
type DeleteTask struct {
	tx       *gorm.DB
	taskID   uint64
	taskRepo repository.TaskRepositoryInterface
	userRepo repository.UserRepositoryInterface
}

// NewDeleteTask constructor
func NewDeleteTask(tx *gorm.DB, taskID uint64, taskRepo repository.TaskRepositoryInterface, userRepo repository.UserRepositoryInterface) *DeleteTask {
	return &DeleteTask{
		tx:       tx,
		taskID:   taskID,
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

// Execute service main process
func (t *DeleteTask) Execute(auth *middleware.Auth) error {
	// get user id from token
	user, err := t.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return err
	}

	// delete task
	task := db.Task{ID: t.taskID, MemberID: user.ID}
	err = t.taskRepo.DeleteTaskWhereMemberID(&task)
	if err != nil {
		return err
	}
	return nil
}
