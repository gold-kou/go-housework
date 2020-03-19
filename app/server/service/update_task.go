package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// UpdateTaskServiceInterface is a service interface of updateTask
type UpdateTaskServiceInterface interface {
	Execute(auth *middleware.Auth) (*db.Task, *db.Family, *db.User, error)
}

// UpdateTask struct
type UpdateTask struct {
	tx               *gorm.DB
	reqUpdateTask    *schemamodel.RequestUpdateTask
	taskRepo         repository.TaskRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewUpdateTask constructor
func NewUpdateTask(tx *gorm.DB, reqUpdateTask *schemamodel.RequestUpdateTask,
	taskRepo repository.TaskRepositoryInterface, userRepo repository.UserRepositoryInterface,
	familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *UpdateTask {
	return &UpdateTask{
		tx:               tx,
		reqUpdateTask:    reqUpdateTask,
		taskRepo:         taskRepo,
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (t *UpdateTask) Execute(auth *middleware.Auth) (*db.Task, *db.Family, *db.User, error) {
	// get user id from token
	user, err := t.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}

	// get family
	familyMember, err := t.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)
	if err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}
	family, err := t.familyRepo.ShowFamily(familyMember.FamilyID)
	if err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}

	// get user id from request parmeter
	reqUser, err := t.userRepo.GetUserWhereUsername(t.reqUpdateTask.Task.MemberName)
	if err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}

	// update task
	dbTask := &db.Task{ID: uint64(t.reqUpdateTask.Task.TaskId), Name: t.reqUpdateTask.Task.TaskName, MemberID: reqUser.ID, FamilyID: family.ID, Status: t.reqUpdateTask.Task.Status, Date: t.reqUpdateTask.Task.Date}
	if err := t.taskRepo.UpdateTask(dbTask); err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}
	return dbTask, family, reqUser, nil
}
