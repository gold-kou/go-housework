package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// ListTasksServiceInterface is a service interface of listTasks
type ListTasksServiceInterface interface {
	Execute(auth *middleware.Auth) ([]*db.Task, *db.Family, []*db.User, error)
}

// ListTasks struct
type ListTasks struct {
	tx               *gorm.DB
	targetDate       string
	taskRepo         repository.TaskRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewListTasks constructor
func NewListTasks(tx *gorm.DB, targetDate string,
	taskRepo repository.TaskRepositoryInterface, userRepo repository.UserRepositoryInterface,
	familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *ListTasks {
	return &ListTasks{
		tx:               tx,
		targetDate:       targetDate,
		taskRepo:         taskRepo,
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (t *ListTasks) Execute(auth *middleware.Auth) ([]*db.Task, *db.Family, []*db.User, error) {
	// get user id from token
	user, err := t.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return nil, &db.Family{}, nil, err
	}

	// get family
	familyMember, err := t.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)
	if err != nil {
		return nil, &db.Family{}, nil, err
	}
	family, err := t.familyRepo.ShowFamily(familyMember.FamilyID)
	if err != nil {
		return nil, &db.Family{}, nil, err
	}

	// select task
	dbTasks, err := t.taskRepo.SelectTaskWhereFamilyIDDate(familyMember.FamilyID, t.targetDate)
	if err != nil {
		return nil, &db.Family{}, nil, err
	}

	// get user name being in charge of the task
	var userIDs []uint64
	for _, t := range dbTasks {
		userIDs = append(userIDs, t.MemberID)
	}
	taskUsers, err := t.userRepo.GetUsersWhereUserIDs(userIDs)
	if err != nil {
		return nil, &db.Family{}, nil, err
	}

	return dbTasks, family, taskUsers, nil
}
