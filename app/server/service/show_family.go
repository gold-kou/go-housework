package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// ShowFamilyServiceInterface is a service interface of showFamily
type ShowFamilyServiceInterface interface {
	Execute() (*db.Family, error)
}

// ShowFamily struct
type ShowFamily struct {
	tx               *gorm.DB
	familyRepo       repository.FamilyRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepository
}

// NewShowFamily constructor
func NewShowFamily(tx *gorm.DB, familyRepo repository.FamilyRepositoryInterface,
	userRepo repository.UserRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepository) *ShowFamily {
	return &ShowFamily{
		tx:               tx,
		familyRepo:       familyRepo,
		userRepo:         userRepo,
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
