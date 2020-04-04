package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// ShowFamilyServiceInterface is a service interface of showFamily
type ShowFamilyServiceInterface interface {
	Execute(*middleware.Auth) (*db.Family, error)
}

// ShowFamily struct
type ShowFamily struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewShowFamily constructor
func NewShowFamily(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *ShowFamily {
	return &ShowFamily{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (f *ShowFamily) Execute(auth *middleware.Auth) (*db.Family, error) {
	// get user id
	user, err := f.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.Family{}, err
	}

	// get family id
	mf, err := f.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)

	// show family
	dbFamily, err := f.familyRepo.ShowFamily(mf.FamilyID)
	if err != nil {
		return &db.Family{}, err
	}

	return dbFamily, nil
}
