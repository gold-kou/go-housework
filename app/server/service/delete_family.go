package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// DeleteFamilyServiceInterface is a service interface of deleteFamily
type DeleteFamilyServiceInterface interface {
	Execute() (*db.Family, error)
}

// DeleteFamily struct
type DeleteFamily struct {
	tx               *gorm.DB
	familyRepo       repository.FamilyRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepository
}

// NewDeleteFamily constructor
func NewDeleteFamily(tx *gorm.DB, familyRepo repository.FamilyRepositoryInterface,
	userRepo repository.UserRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepository) *DeleteFamily {
	return &DeleteFamily{
		tx:               tx,
		familyRepo:       familyRepo,
		userRepo:         userRepo,
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
