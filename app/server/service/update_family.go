package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// UpdateFamilyServiceInterface is a service interface of updateFamily
type UpdateFamilyServiceInterface interface {
	Execute(*middleware.Auth, *schemamodel.RequestUpdateFamily) (*db.Family, error)
}

// UpdateFamily struct
type UpdateFamily struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewUpdateFamily constructor
func NewUpdateFamily(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *UpdateFamily {
	return &UpdateFamily{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (f *UpdateFamily) Execute(auth *middleware.Auth, updateFamily *schemamodel.RequestUpdateFamily) (*db.Family, error) {
	// get user id
	user, err := f.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.Family{}, err
	}

	// get family id
	mf, err := f.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)

	// data
	dbFamily := &db.Family{ID: mf.FamilyID, Name: updateFamily.FamilyName}

	// update family
	if err := f.familyRepo.UpdateFamily(dbFamily); err != nil {
		return &db.Family{}, err
	}

	return dbFamily, nil
}
