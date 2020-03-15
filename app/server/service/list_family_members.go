package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// ListFamilyMembersServiceInterface is a service interface of listFamilyMembers
type ListFamilyMembersServiceInterface interface {
	Execute() (*db.Family, []*db.User, error)
}

// ListFamilyMembers struct
type ListFamilyMembers struct {
	tx               *gorm.DB
	familyRepo       repository.FamilyRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepository
}

// NewListFamilyMembers constructor
func NewListFamilyMembers(tx *gorm.DB, familyRepo repository.FamilyRepositoryInterface,
	userRepo repository.UserRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepository) *ListFamilyMembers {
	return &ListFamilyMembers{
		tx:               tx,
		familyRepo:       familyRepo,
		userRepo:         userRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (f *ListFamilyMembers) Execute(auth *middleware.Auth) (*db.Family, []*db.User, error) {
	// get user id from token
	user, err := f.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.Family{}, nil, err
	}

	// get family_id by user_id
	mf, err := f.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)
	if err != nil {
		return &db.Family{}, nil, err
	}

	// get family
	dbFamily, err := f.familyRepo.ShowFamily(mf.FamilyID)
	if err != nil {
		return &db.Family{}, nil, err
	}

	// get all user_id which belongs to the family
	mfs, err := f.memberFamilyRepo.ListMemberFamiliesWhereFamilyID(mf.FamilyID)
	if err != nil {
		return &db.Family{}, nil, err
	}

	// get all users
	var users []*db.User
	for _, memberFamily := range mfs {
		u, err := f.userRepo.GetUserWhereUserID(memberFamily.MemberID)
		if err != nil {
			return &db.Family{}, nil, err
		}
		users = append(users, u)
	}

	return dbFamily, users, nil
}
