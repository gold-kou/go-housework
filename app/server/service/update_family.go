package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// UpdateFamilyServiceInterface is a service interface of updateFamily
type UpdateFamilyServiceInterface interface {
	Execute() (*db.Family, error)
}

// UpdateFamily struct
type UpdateFamily struct {
	tx               *gorm.DB
	updateFamily     *schemamodel.RequestUpdateFamily
	familyRepo       repository.FamilyRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepository
}

// NewUpdateFamily constructor
func NewUpdateFamily(tx *gorm.DB, updateFamily *schemamodel.RequestUpdateFamily,
	familyRepo repository.FamilyRepositoryInterface, userRepo repository.UserRepositoryInterface,
	memberFamilyRepo repository.MemberFamilyRepository) *UpdateFamily {
	return &UpdateFamily{
		tx:               tx,
		updateFamily:     updateFamily,
		familyRepo:       familyRepo,
		userRepo:         userRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (f *UpdateFamily) Execute(auth *middleware.Auth) (*db.Family, error) {
	// get user id
	user, err := f.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.Family{}, err
	}

	// get family id
	mf, err := f.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)

	// data
	dbFamily := &db.Family{ID: mf.FamilyID, Name: f.updateFamily.FamilyName}

	// update family
	if err := f.familyRepo.UpdateFamily(dbFamily); err != nil {
		return &db.Family{}, err
	}

	return dbFamily, nil
}
