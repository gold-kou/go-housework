package service

import (
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserServiceInterface is a service interface of createUser
type CreateUserServiceInterface interface {
	Execute(createUser *schemamodel.RequestCreateUser) (*db.User, error)
}

// CreateUser struct
type CreateUser struct {
	userRepo repository.UserRepositoryInterface
}

// NewCreateUser constructor
func NewCreateUser(userRepo repository.UserRepositoryInterface) *CreateUser {
	return &CreateUser{userRepo: userRepo}
}

func (u *CreateUser) Execute(createUser *schemamodel.RequestCreateUser) (*db.User, error) {
	// data
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Password), bcrypt.DefaultCost)
	dbUser := &db.User{Name: createUser.UserName, Email: createUser.Email, Password: string(hashedPassword)}

	// insert user
	if err = u.userRepo.InsertUser(dbUser); err != nil {
		return &db.User{}, err
	}
	return dbUser, nil
}
