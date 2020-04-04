package service

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Execute(t *testing.T) {
	type args struct {
		userName string
		password string
	}
	tests := []struct {
		name             string
		mockUserRepoFunc func(*repository.MockUserRepositoryInterface)
		args             args
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "success",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().
					GetUserWhereUsername(common.TestUserName).
					Return(&db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestHashedPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()}, nil)
			},
			args:    args{userName: common.TestUserName, password: common.TestPassword},
			wantErr: false,
		},
		{
			name: "fail(not found user)",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().
					GetUserWhereUsername(common.TestUserName).
					Return(&db.User{}, errors.New("not found user"))
			},
			args:       args{userName: common.TestUserName, password: common.TestWrongPassword},
			wantErr:    true,
			wantErrMsg: "not found user",
		},
		{
			name: "fail(wrong password)",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().
					GetUserWhereUsername(common.TestUserName).
					Return(&db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestHashedPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()}, nil)
			},
			args:       args{userName: common.TestUserName, password: common.TestWrongPassword},
			wantErr:    true,
			wantErrMsg: "not correct password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// connect dummy db
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				// mock repository
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				userRepo := repository.NewMockUserRepositoryInterface(ctrl)
				tt.mockUserRepoFunc(userRepo)

				// run target method
				_, err := NewLogin(userRepo).Execute(tt.args.userName, tt.args.password)

				// assert
				assert := assert.New(t)
				if tt.wantErr {
					assert.Error(err)
					assert.EqualError(err, tt.wantErrMsg)
				} else {
					assert.NoError(err)
					// MEMO: JWTの中身は毎実行する度に変わるのでテスト省略
				}
			})
		})
	}
}
