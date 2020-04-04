package service

import (
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// DeleteFamilyServiceInterface is a service interface of deleteFamily
type DeleteFamilyServiceInterface interface {
	Execute(*middleware.Auth) error
}

// DeleteFamily struct
type DeleteFamily struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewDeleteFamily constructor
func NewDeleteFamily(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *DeleteFamily {
	return &DeleteFamily{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (f *DeleteFamily) Execute(auth *middleware.Auth) error {
	// get user id
	user, err := f.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return err
	}

	// get family id
	mf, err := f.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)

	// delete member_family
	if err = f.memberFamilyRepo.DeleteMemberFamily(mf.FamilyID); err != nil {
		return err
	}

	// delete family
	if err := f.familyRepo.DeleteFamily(mf.FamilyID); err != nil {
		return err
	}

	return nil
}
