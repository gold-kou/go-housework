package service

import (
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// DeleteFamilyMemberServiceInterface is a service interface of deleteFamilyMember
type DeleteFamilyMemberServiceInterface interface {
	Execute(*middleware.Auth, uint64) error
}

// DeleteFamilyMember struct
type DeleteFamilyMember struct {
	userRepo         repository.UserRepositoryInterface
	familyRepo       repository.FamilyRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewDeleteFamilyMember constructor
func NewDeleteFamilyMember(userRepo repository.UserRepositoryInterface, familyRepo repository.FamilyRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepositoryInterface) *DeleteFamilyMember {
	return &DeleteFamilyMember{
		userRepo:         userRepo,
		familyRepo:       familyRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (mf *DeleteFamilyMember) Execute(auth *middleware.Auth, memberID uint64) error {
	// delete member_family
	if err := mf.memberFamilyRepo.DeleteMemberFromFamily(memberID); err != nil {
		return err
	}
	return nil
}
