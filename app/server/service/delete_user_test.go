package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_Execute(t *testing.T) {
	type args struct {
		auth *model.Auth
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
				r.EXPECT().DeleteUserWhereUsername(common.TestUserName).Return(nil)
			},
			args:    args{&model.Auth{UserName: common.TestUserName}},
			wantErr: false,
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
				err := NewDeleteUser(userRepo).Execute(tt.args.auth)

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
