package service

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask_Execute(t *testing.T) {
	type args struct {
		auth       *model.Auth
		createTask *schemamodel.RequestCreateTask
	}
	tests := []struct {
		name                     string
		mockUserRepoFunc         func(*repository.MockUserRepositoryInterface)
		mockFamilyRepoFunc       func(*repository.MockFamilyRepositoryInterface)
		mockMemberFamilyRepoFunc func(*repository.MockMemberFamilyRepositoryInterface)
		mockTaskRepoFunc         func(*repository.MockTaskRepositoryInterface)
		args                     args
		want1                    *db.User
		want2                    *db.Family
		want3                    *db.Task
		wantErr                  bool
		wantErrMsg               string
	}{
		{
			name: "success",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().GetUserWhereUsername(common.TestUserName).Return(&db.User{ID: common.TestUserID, Name: common.TestUserName}, nil) // for token user
				r.EXPECT().GetUserWhereUsername(common.TestUserName).Return(&db.User{ID: common.TestUserID, Name: common.TestUserName}, nil) // for request parameter user
			},
			mockFamilyRepoFunc: func(r *repository.MockFamilyRepositoryInterface) {
				r.EXPECT().ShowFamily(common.TestFamilyID).Return(&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, nil)
			},
			mockMemberFamilyRepoFunc: func(r *repository.MockMemberFamilyRepositoryInterface) {
				r.EXPECT().GetMemberFamilyWhereMemberID(common.TestUserID).Return(&db.MemberFamily{MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Role: common.TestRoleHead}, nil)
			},
			mockTaskRepoFunc: func(r *repository.MockTaskRepositoryInterface) {
				r.EXPECT().InsertTask(&db.Task{Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate}).Return(nil)
			},
			args: args{
				auth:       &model.Auth{UserName: common.TestUserName},
				createTask: &schemamodel.RequestCreateTask{TaskName: common.TestTaskName1, MemberName: common.TestUserName, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate},
			},
			want1:   &db.User{ID: common.TestUserID, Name: common.TestUserName},
			want2:   &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName},
			want3:   &db.Task{Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate}, // MEMO: 参照渡しなのでFamilyIDが連携されるべきだがmockではそこまで実現できないためID省略
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
				familyRepo := repository.NewMockFamilyRepositoryInterface(ctrl)
				tt.mockFamilyRepoFunc(familyRepo)
				memberFamilyRepo := repository.NewMockMemberFamilyRepositoryInterface(ctrl)
				tt.mockMemberFamilyRepoFunc(memberFamilyRepo)
				taskRepo := repository.NewMockTaskRepositoryInterface(ctrl)
				tt.mockTaskRepoFunc(taskRepo)

				// run target method
				got1, got2, got3, err := NewCreateTask(userRepo, familyRepo, memberFamilyRepo, taskRepo).Execute(tt.args.auth, tt.args.createTask)

				// assert
				assert := assert.New(t)
				if tt.wantErr {
					assert.Error(err)
					assert.EqualError(err, tt.wantErrMsg)
				} else {
					assert.NoError(err)
					assert.Equal(tt.want1, got1)
					assert.Equal(tt.want2, got2)
					assert.Equal(tt.want3, got3)
				}
			})
		})
	}
}
