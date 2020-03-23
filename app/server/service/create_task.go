package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// CreateTaskServiceInterface is a service interface of createTask
type CreateTaskServiceInterface interface {
	Execute(auth *middleware.Auth) (*db.Task, *db.Family, *db.User, error)
}

// CreateTask struct
type CreateTask struct {
	tx               *gorm.DB
	reqCreateTask    *schemamodel.RequestCreateTask
	taskRepo         repository.TaskRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewCreateTask constructor
func NewCreateTask(tx *gorm.DB, reqCreateTask *schemamodel.RequestCreateTask,
	taskRepo repository.TaskRepositoryInterface, userRepo repository.UserRepositoryInterface,
	familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *CreateTask {
	return &CreateTask{
		tx:               tx,
		reqCreateTask:    reqCreateTask,
		taskRepo:         taskRepo,
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (t *CreateTask) Execute(auth *middleware.Auth) (*db.Task, *db.Family, *db.User, error) {
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
	reqUser, err := t.userRepo.GetUserWhereUsername(t.reqCreateTask.MemberName)
	if err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}

	// insert task
	dbTask := &db.Task{Name: t.reqCreateTask.TaskName, MemberID: reqUser.ID, FamilyID: family.ID, Status: t.reqCreateTask.Status, Date: t.reqCreateTask.Date}
	if err := t.taskRepo.InsertTask(dbTask); err != nil {
		return &db.Task{}, &db.Family{}, &db.User{}, err
	}
	return dbTask, family, reqUser, nil
}
