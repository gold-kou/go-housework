package service

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// CreateFamilyServiceInterface is a service interface of createFamily
type CreateFamilyServiceInterface interface {
	Execute() (*db.Family, error)
}

// CreateFamily struct
type CreateFamily struct {
	tx               *gorm.DB
	createFamily     *schemamodel.RequestCreateFamily
	familyRepo       repository.FamilyRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepository
}

// NewCreateFamily constructor
func NewCreateFamily(tx *gorm.DB, createFamily *schemamodel.RequestCreateFamily,
	familyRepo repository.FamilyRepositoryInterface, userRepo repository.UserRepositoryInterface,
	memberFamilyRepo repository.MemberFamilyRepository) *CreateFamily {
	return &CreateFamily{
		tx:               tx,
		createFamily:     createFamily,
		familyRepo:       familyRepo,
		userRepo:         userRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (f *CreateFamily) Execute(auth *middleware.Auth) (*db.Family, error) {
	// get user id
	user, err := f.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.Family{}, err
	}

	// check if the user already belongs to any families
	mf, err := f.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)
	if mf != nil {
		return &db.Family{}, common.NewBadRequestError("the user has already belonged to another family")
	}

	// data
	dbFamily := &db.Family{Name: f.createFamily.FamilyName}

	// insert family
	if err := f.familyRepo.InsertFamily(dbFamily); err != nil {
		return &db.Family{}, err
	}

	// insert member_family
	memberFamily := db.MemberFamily{FamilyID: dbFamily.ID, MemberID: user.ID, Role: db.FamilyRoleHead}
	if err = f.memberFamilyRepo.InsertMemberFamily(&memberFamily); err != nil {
		return &db.Family{}, err
	}

	return dbFamily, nil
}
