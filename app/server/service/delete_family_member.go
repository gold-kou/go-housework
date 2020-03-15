package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
)

// DeleteFamilyMemberServiceInterface is a service interface of deleteFamilyMember
type DeleteFamilyMemberServiceInterface interface {
	Execute() (*db.Family, error)
}

// DeleteFamilyMember struct
type DeleteFamilyMember struct {
	tx               *gorm.DB
	memberID         uint64
	familyRepo       repository.FamilyRepositoryInterface
	userRepo         repository.UserRepositoryInterface
	memberFamilyRepo repository.MemberFamilyRepository
}

// NewDeleteFamilyMember constructor
func NewDeleteFamilyMember(tx *gorm.DB, memberID uint64, familyRepo repository.FamilyRepositoryInterface,
	userRepo repository.UserRepositoryInterface, memberFamilyRepo repository.MemberFamilyRepository) *DeleteFamilyMember {
	return &DeleteFamilyMember{
		tx:               tx,
		memberID:         memberID,
		familyRepo:       familyRepo,
		userRepo:         userRepo,
		memberFamilyRepo: memberFamilyRepo,
	}
}

// Execute service main process
func (mf *DeleteFamilyMember) Execute(auth *middleware.Auth) error {
	// delete member_family
	if err := mf.memberFamilyRepo.DeleteMemberFromFamily(mf.memberID); err != nil {
		return err
	}

	return nil
}
