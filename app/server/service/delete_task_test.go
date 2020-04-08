package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTask_Execute(t *testing.T) {
	type args struct {
		auth   *model.Auth
		taskID uint64
	}
	tests := []struct {
		name             string
		mockUserRepoFunc func(*repository.MockUserRepositoryInterface)
		mockTaskRepoFunc func(*repository.MockTaskRepositoryInterface)
		args             args
		wantErr          bool
		wantErrMsg       string
	}{
		{
			name: "success",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().GetUserWhereUsername(common.TestUserName).Return(&db.User{ID: common.TestUserID, Name: common.TestUserName}, nil)
			},
			mockTaskRepoFunc: func(r *repository.MockTaskRepositoryInterface) {
				r.EXPECT().DeleteTaskWhereMemberID(&db.Task{ID: common.TestTaskID1, MemberID: common.TestUserID}).Return(nil)
			},
			args: args{
				auth:   &model.Auth{UserName: common.TestUserName},
				taskID: common.TestTaskID1,
			},
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
				taskRepo := repository.NewMockTaskRepositoryInterface(ctrl)
				tt.mockTaskRepoFunc(taskRepo)

				// run target method
				err := NewDeleteTask(userRepo, taskRepo).Execute(tt.args.auth, tt.args.taskID)

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
