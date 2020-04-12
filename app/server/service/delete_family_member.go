package service

import (
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/server/repository"
)

// DeleteFamilyMemberServiceInterface is a service interface of deleteFamilyMember
type DeleteFamilyMemberServiceInterface interface {
	Execute(*model.Auth, uint64) error
}

// DeleteFamilyMember struct
type DeleteFamilyMember struct {
	memberFamilyRepo repository.MemberFamilyRepositoryInterface
}

// NewDeleteFamilyMember constructor
func NewDeleteFamilyMember(memberFamilyRepo repository.MemberFamilyRepositoryInterface) *DeleteFamilyMember {
	return &DeleteFamilyMember{
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (mf *DeleteFamilyMember) Execute(auth *model.Auth, memberID uint64) error {
	// delete member_family
	if err := mf.memberFamilyRepo.DeleteMemberFromFamily(memberID); err != nil {
		return err
	}
	return nil
}
