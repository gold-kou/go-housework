package repository

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/jinzhu/gorm"
)

// UserRepositoryInterface is a repository interface of user
type UserRepositoryInterface interface {
	InsertUser(*db.User) error
	GetUserWhereUsername(string) (*db.User, error)
	DeleteUserWhereUsername(string) error
}

// UserRepository is a repository of user
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a pointer of UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// InsertUser insert user
func (u *UserRepository) InsertUser(user *db.User) error {
	if err := u.db.Create(&user).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}

// GetUserWhereUsername returns user. returns nil if there is no record.
func (u *UserRepository) GetUserWhereUsername(userName string) (*db.User, error) {
	var user db.User
	if err := u.db.Where("name = ?", userName).Find(&user).Error; err != nil {
		// no user record
		if gorm.IsRecordNotFoundError(err) {
			return &db.User{}, common.NewBadRequestError("not found user")
		}
		// unexpected error
		return &db.User{}, common.NewInternalServerError(err.Error())
	}
	return &user, nil
}

// DeleteUserWhereUsername delete user
func (u *UserRepository) DeleteUserWhereUsername(userName string) error {
	var user db.User
	if err := u.db.Where("name = ?", userName).Delete(&user).Error; err != nil {
		return common.NewInternalServerError(err.Error())
	}
	return nil
}
