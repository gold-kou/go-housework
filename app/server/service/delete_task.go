package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// DeleteTaskServiceInterface is a service interface of deleteTask
type DeleteTaskServiceInterface interface {
	Execute(*middleware.Auth, uint64) error
}

// DeleteTask struct
type DeleteTask struct {
	userRepo repository.UserRepositoryInterface
	taskRepo repository.TaskRepositoryInterface
}

// NewDeleteTask constructor
func NewDeleteTask(userRepo repository.UserRepositoryInterface, taskRepo repository.TaskRepositoryInterface) *DeleteTask {
	return &DeleteTask{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

// Execute service main process
func (t *DeleteTask) Execute(auth *middleware.Auth, taskID uint64) error {
	// get user id from token
	user, err := t.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return err
	}

	// delete task
	task := db.Task{ID: taskID, MemberID: user.ID}
	err = t.taskRepo.DeleteTaskWhereMemberID(&task)
	if err != nil {
		return err
	}
	return nil
}
