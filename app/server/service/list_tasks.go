package service

import (
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// ListTasksServiceInterface is a service interface of listTasks
type ListTasksServiceInterface interface {
	Execute(*model.Auth, string) ([]*db.Task, *db.Family, []*db.User, error)
}

// ListTasks struct
type ListTasks struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
	taskRepo         repository.TaskRepositoryInterface
}

// NewListTasks constructor
func NewListTasks(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface, taskRepo repository.TaskRepositoryInterface) *ListTasks {
	return &ListTasks{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
		taskRepo:         taskRepo,
	}
}

// Execute service main process
func (t *ListTasks) Execute(auth *model.Auth, targetDate string) ([]*db.Task, *db.Family, []*db.User, error) {
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
	dbTasks, err := t.taskRepo.SelectTaskWhereFamilyIDDate(familyMember.FamilyID, targetDate)
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
