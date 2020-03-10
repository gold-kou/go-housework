package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser_Execute(t *testing.T) {
	type args struct {
		user *db.User
	}
	tests := []struct {
		name             string
		mockUserRepoFunc func(*repository.MockUserRepositoryInterface)
		args             args
		wantErr          bool
		wantErrMsg       string
	}{
		/*
			// TODO because of bcrypt randam
			{
				name: "success",
				mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
					r.EXPECT().
						InsertUser(&db.User{}).
						Return(nil)
				},
				args: args{user: &db.User{Name: common.TestUserName, Email: common.TestEmail, Password: common.TestPassword}},
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
				reqCreateUser := &schemamodel.RequestCreateUser{UserName: common.TestUserName, Email: common.TestEmail, Password: common.TestPassword}
				_, err := NewCreateUser(db, reqCreateUser, userRepo).Execute()

				// assert
				assert := assert.New(t)
				if tt.wantErr {
					assert.Error(err)
					assert.EqualError(err, tt.wantErrMsg)
				} else {
					assert.NoError(err)
					// assert.Equal(db.User{ID:common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword, CreatedAt:common.GetTestTime(), UpdatedAt:common.GetTestTime()}, got)
				}
			})
		})
	}
}
