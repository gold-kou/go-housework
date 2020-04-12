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

func TestListTasks_Execute(t *testing.T) {
	dbExpectedTasks := make([]*db.Task, 0)
	dbExpectedTasks = append(dbExpectedTasks, &db.Task{ID: common.TestTaskID1, Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate})

	dbExpectedUsers := make([]*db.User, 0)
	dbExpectedUsers = append(dbExpectedUsers, &db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword})

	type args struct {
		auth       *model.Auth
		targetDate string
	}
	tests := []struct {
		name                     string
		mockUserRepoFunc         func(*repository.MockUserRepositoryInterface)
		mockFamilyRepoFunc       func(*repository.MockFamilyRepositoryInterface)
		mockMemberFamilyRepoFunc func(*repository.MockMemberFamilyRepositoryInterface)
		mockTaskRepoFunc         func(*repository.MockTaskRepositoryInterface)
		args                     args
		want1                    []*db.Task
		want2                    *db.Family
		want3                    []*db.User
		wantErr                  bool
		wantErrMsg               string
	}{
		{
			name: "success",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().GetUserWhereUsername(common.TestUserName).Return(&db.User{ID: common.TestUserID, Name: common.TestUserName}, nil)
				r.EXPECT().GetUsersWhereUserIDs([]uint64{common.TestUserID}).Return(dbExpectedUsers, nil)
			},
			mockFamilyRepoFunc: func(r *repository.MockFamilyRepositoryInterface) {
				r.EXPECT().ShowFamily(common.TestFamilyID).Return(&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, nil)
			},
			mockMemberFamilyRepoFunc: func(r *repository.MockMemberFamilyRepositoryInterface) {
				r.EXPECT().GetMemberFamilyWhereMemberID(common.TestUserID).Return(&db.MemberFamily{MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Role: common.TestRoleHead}, nil)
			},
			mockTaskRepoFunc: func(r *repository.MockTaskRepositoryInterface) {
				r.EXPECT().SelectTaskWhereFamilyIDDate(common.TestFamilyID, common.TestTaskDate).Return(dbExpectedTasks, nil)
			},
			args: args{
				auth:       &model.Auth{UserName: common.TestUserName},
				targetDate: common.TestTaskDate,
			},
			want1:   dbExpectedTasks,
			want2:   &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName},
			want3:   dbExpectedUsers,
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
				got1, got2, got3, err := NewListTasks(userRepo, familyRepo, memberFamilyRepo, taskRepo).Execute(tt.args.auth, tt.args.targetDate)

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
