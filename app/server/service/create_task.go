package service

import (
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// CreateTaskServiceInterface is a service interface of createTask
type CreateTaskServiceInterface interface {
	Execute(*model.Auth, *schemamodel.RequestCreateTask) (*db.User, *db.Family, *db.Task, error)
}

// CreateTask struct
type CreateTask struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
	taskRepo         repository.TaskRepositoryInterface
}

// NewCreateTask constructor
func NewCreateTask(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface,
	memberFamilyRepo repository.MemberFamilyRepositoryInterface, taskRepo repository.TaskRepositoryInterface) *CreateTask {
	return &CreateTask{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
		taskRepo:         taskRepo,
	}
}

// Execute service main process
func (t *CreateTask) Execute(auth *model.Auth, reqCreateTask *schemamodel.RequestCreateTask) (*db.User, *db.Family, *db.Task, error) {
	// get user id from token
	user, err := t.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.User{}, &db.Family{}, &db.Task{}, err
	}

	// get family
	familyMember, err := t.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)
	if err != nil {
		return &db.User{}, &db.Family{}, &db.Task{}, err
	}
	family, err := t.familyRepo.ShowFamily(familyMember.FamilyID)
	if err != nil {
		return &db.User{}, &db.Family{}, &db.Task{}, err
	}

	// get user id from request parmeter
	reqUser, err := t.userRepo.GetUserWhereUsername(reqCreateTask.MemberName)
	if err != nil {
		return &db.User{}, &db.Family{}, &db.Task{}, err
	}

	// insert task
	dbTask := &db.Task{Name: reqCreateTask.TaskName, MemberID: reqUser.ID, FamilyID: family.ID, Status: reqCreateTask.Status, Date: reqCreateTask.Date}
	if err := t.taskRepo.InsertTask(dbTask); err != nil {
		return &db.User{}, &db.Family{}, &db.Task{}, err
	}
	return reqUser, family, dbTask, nil
}
