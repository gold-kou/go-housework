package service

import (
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// LoginServiceInterface is a service interface of login
type LoginServiceInterface interface {
	Execute() (string, error)
}

// Login struct
type Login struct {
	tx       *gorm.DB
	userName string
	password string
	userRepo repository.UserRepositoryInterface
}

// NewLogin constructor
func NewLogin(tx *gorm.DB, userName string, password string, userRepo repository.UserRepositoryInterface) *Login {
	return &Login{
		tx:       tx,
		userName: userName,
		password: password,
		userRepo: userRepo,
	}
}

func (l *Login) Execute() (string, error) {
	// ユーザ名からユーザ情報の取得
	user, err := l.userRepo.GetUserWhereUsername(l.userName)
	if err != nil {
		return "", err
	}

	// パスワードの検証(DBに格納されたハッシュ値とクエリパラメータで渡された値を使った比較)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(l.password))
	if err != nil {
		return "", common.NewAuthorizationError("not correct password")
	}

	// JWT生成
	tokenString, err := middleware.GenerateToken(l.userName)
	if err != nil {
		return "", common.NewInternalServerError(err.Error())
	}
	return tokenString, nil
}
