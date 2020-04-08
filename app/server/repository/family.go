package repository

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/jinzhu/gorm"
)

// FamilyRepositoryInterface is a repository interface of family
type FamilyRepositoryInterface interface {
	InsertFamily(*db.Family) error
	UpdateFamily(*db.Family) error
	DeleteFamily(uint64) error
	ShowFamily(uint64) (*db.Family, error)
}

// FamilyRepository is a repository of family
type FamilyRepository struct {
	db *gorm.DB
}

// NewFamilyRepository creates a pointer of FamilyRepository
func NewFamilyRepository(db *gorm.DB) *FamilyRepository {
	return &FamilyRepository{
		db: db,
	}
}

// InsertFamily insert family
func (f *FamilyRepository) InsertFamily(family *db.Family) error {
	if err := f.db.Create(&family).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// UpdateFamily update family
func (f *FamilyRepository) UpdateFamily(family *db.Family) error {
	if err := f.db.Save(family).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// DeleteFamily delete family
func (f *FamilyRepository) DeleteFamily(familyID uint64) error {
	var family db.Family
	if err := f.db.Where("id = ?", familyID).Delete(&family).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// ShowFamily show family
func (f *FamilyRepository) ShowFamily(familyID uint64) (*db.Family, error) {
	var family db.Family
	if err := f.db.Where("id = ?", familyID).Find(&family).Error; err != nil {
		return &db.Family{}, common.NewInternalServerError(err.Error())
	}
	return &family, nil
}
