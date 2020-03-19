package repository

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/jinzhu/gorm"
)

// MemberFamilyRepositoryInterface is a repository interface of memberFamily
type MemberFamilyRepositoryInterface interface {
	InsertMemberFamily(*db.MemberFamily) error
	DeleteMemberFamily(uint64) error
	DeleteMemberFromFamily(uint64) error
	GetMemberFamilyWhereMemberID(uint64) (*db.MemberFamily, error)
	ListMemberFamiliesWhereFamilyID(uint64) ([]*db.MemberFamily, error)
}

// MemberFamilyRepository is a repository of memberFamily
type MemberFamilyRepository struct {
	db *gorm.DB
}

// NewMemberFamilyRepository creates a pointer of MemberFamilyRepository
func NewMemberFamilyRepository(db *gorm.DB) *MemberFamilyRepository {
	return &MemberFamilyRepository{
		db: db,
	}
}

// InsertMemberFamily insert memberFamily
func (mf *MemberFamilyRepository) InsertMemberFamily(memberFamily *db.MemberFamily) error {
	if err := mf.db.Create(&memberFamily).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// DeleteMemberFamily delete memberFamily
func (mf *MemberFamilyRepository) DeleteMemberFamily(familyID uint64) error {
	var memberFamily db.MemberFamily
	if err := mf.db.Where("family_id = ?", familyID).Delete(&memberFamily).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// DeleteMemberFromFamily delete member from memberFamily
func (mf *MemberFamilyRepository) DeleteMemberFromFamily(memberID uint64) error {
	var memberFamily db.MemberFamily
	if err := mf.db.Where("member_id = ?", memberID).Delete(&memberFamily).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// GetMemberFamilyWhereMemberID get memberFamily
func (mf *MemberFamilyRepository) GetMemberFamilyWhereMemberID(memberID uint64) (*db.MemberFamily, error) {
	var memberFamily db.MemberFamily
	if err := mf.db.Where("member_id = ?", memberID).Find(&memberFamily).Error; err != nil {
		// empty is not error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		// unexpected error
		return &memberFamily, common.NewInternalServerError(err.Error())
	}
	// not empty
	return &memberFamily, nil
}

// ListMemberFamiliesWhereFamilyID get memberFamilies
func (mf *MemberFamilyRepository) ListMemberFamiliesWhereFamilyID(familyID uint64) ([]*db.MemberFamily, error) {
	var memberFamilies []*db.MemberFamily
	if err := mf.db.Where("family_id = ?", familyID).Order("member_id").Find(&memberFamilies).Error; err != nil {
		return memberFamilies, common.NewInternalServerError(err.Error())
	}
	return memberFamilies, nil
}
