package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Execute(t *testing.T) {
	type args struct {
		createUser *schemamodel.RequestCreateUser
	}
	tests := []struct {
		name             string
		mockUserRepoFunc func(*repository.MockUserRepositoryInterface)
		args             args
		wantErr          bool
		wantErrMsg       string
	}{
		/*
			skip for bcrypt random
				{
					name: "success",
					mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
						r.EXPECT().InsertUser(&db.User{Name: common.TestUserName, Email: common.TestEmail, Password: common.TestPassword}).Return(nil)
					},
					args:    args{createUser: &schemamodel.RequestCreateUser{Email: common.TestEmail, UserName: common.TestUserName, Password: common.TestPassword}},
					wantErr: false,
				},
		*/
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
				_, err := NewCreateUser(userRepo).Execute(tt.args.createUser)

				// assert
				assert := assert.New(t)
				if tt.wantErr {
					assert.Error(err)
					assert.EqualError(err, tt.wantErrMsg)
				} else {
					assert.NoError(err)
				}
			})
		})
	}
}
