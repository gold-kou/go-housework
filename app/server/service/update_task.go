package service

import (
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// UpdateTaskServiceInterface is a service interface of updateTask
type UpdateTaskServiceInterface interface {
	Execute(*model.Auth, *schemamodel.RequestUpdateTask) (*db.Task, *db.Family, *db.User, error)
}

// UpdateTask struct
type UpdateTask struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
	taskRepo         repository.TaskRepositoryInterface
}

// NewUpdateTask constructor
func NewUpdateTask(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface,
	memberFamilyRepo repository.MemberFamilyRepositoryInterface, taskRepo repository.TaskRepositoryInterface) *UpdateTask {
	return &UpdateTask{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
		taskRepo:         taskRepo,
	}
}

// Execute service main process
func (t *UpdateTask) Execute(auth *model.Auth, reqUpdateTask *schemamodel.RequestUpdateTask) (*db.Task, *db.Family, *db.User, error) {
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
	reqUser, err := t.userRepo.GetUserWhereUsername(reqUpdateTask.Task.MemberName)
	if err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}

	// update task
	dbTask := &db.Task{ID: uint64(reqUpdateTask.Task.TaskId), Name: reqUpdateTask.Task.TaskName, MemberID: reqUser.ID, FamilyID: family.ID, Status: reqUpdateTask.Task.Status, Date: reqUpdateTask.Task.Date}
	if err := t.taskRepo.UpdateTask(dbTask); err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}
	return dbTask, family, reqUser, nil
}
