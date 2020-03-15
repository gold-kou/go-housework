package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserServiceInterface is a service interface of createUser
type CreateUserServiceInterface interface {
	Execute() (*db.User, error)
}

// CreateUser struct
type CreateUser struct {
	tx         *gorm.DB
	createUser *schemamodel.RequestCreateUser
	userRepo   repository.UserRepositoryInterface
}

// NewCreateUser constructor
func NewCreateUser(tx *gorm.DB, createUser *schemamodel.RequestCreateUser, userRepo repository.UserRepositoryInterface) *CreateUser {
	return &CreateUser{
		tx:         tx,
		createUser: createUser,
		userRepo:   userRepo,
	}
}

// Execute service main process
func (u *CreateUser) Execute() (*db.User, error) {
	// data
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.createUser.Password), bcrypt.DefaultCost)
	dbUser := &db.User{Name: u.createUser.UserName, Email: u.createUser.Email, Password: string(hashedPassword)}

	// insert user
	if err = u.userRepo.InsertUser(dbUser); err != nil {
		return &db.User{}, err
	}
	return dbUser, nil
}
