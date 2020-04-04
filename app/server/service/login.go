package service

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"golang.org/x/crypto/bcrypt"
)

// LoginServiceInterface is a service interface of login
type LoginServiceInterface interface {
	Execute(string, string) (string, error)
}

// Login struct
type Login struct {
	userRepo repository.UserRepositoryInterface
}

// NewLogin constructor
func NewLogin(userRepo repository.UserRepositoryInterface) *Login {
	return &Login{
		userRepo: userRepo,
	}
}

// Execute service main process
func (l *Login) Execute(userName string, password string) (string, error) {
	// ユーザ名からユーザ情報の取得
	user, err := l.userRepo.GetUserWhereUsername(userName)
	if err != nil {
		return "", err
	}

	// パスワードの検証(DBに格納されたハッシュ値とクエリパラメータで渡された値を使った比較)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", common.NewAuthorizationError("not correct password")
	}

	// JWT生成
	tokenString, err := middleware.GenerateToken(userName)
	if err != nil {
		return "", common.NewInternalServerError(err.Error())
	}
	return tokenString, nil
}
