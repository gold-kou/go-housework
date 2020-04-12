package service

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// RegisterFamilyMemberServiceInterface is a service interface of registerFamilyMember
type RegisterFamilyMemberServiceInterface interface {
	Execute(*model.Auth, *schemamodel.RequestRegisterFamilyMember) (*db.User, *db.Family, error)
}

// RegisterFamilyMember struct
type RegisterFamilyMember struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewRegisterFamilyMember constructor
func NewRegisterFamilyMember(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *RegisterFamilyMember {
	return &RegisterFamilyMember{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (fm *RegisterFamilyMember) Execute(auth *model.Auth, registerFamilyMember *schemamodel.RequestRegisterFamilyMember) (*db.User, *db.Family, error) {
	// get user_id from auth
	user, err := fm.userRepo.GetUserWhereUsername(auth.UserName)
	if err != nil {
		return &db.User{}, &db.Family{}, err
	}

	// get family id of auth user
	mf, err := fm.memberFamilyRepo.GetMemberFamilyWhereMemberID(user.ID)

	// check role
	if mf.Role != db.FamilyRoleHead {
		return &db.User{}, &db.Family{}, common.NewAuthorizationError("this is not head user")
	}

	// get family from family_id
	dbFamily, err := fm.familyRepo.ShowFamily(mf.FamilyID)
	if err != nil {
		return &db.User{}, &db.Family{}, err
	}

	// get user_id from request param
	targetUser, err := fm.userRepo.GetUserWhereUsername(registerFamilyMember.MemberName)
	if err != nil {
		return &db.User{}, &db.Family{}, err
	}

	// insert member_family
	dbMemberFamily := db.MemberFamily{MemberID: targetUser.ID, FamilyID: mf.FamilyID, Role: db.FamilyRoleMember}
	if err = fm.memberFamilyRepo.InsertMemberFamily(&dbMemberFamily); err != nil {
		return &db.User{}, &db.Family{}, err
	}

	return targetUser, dbFamily, nil
}
