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

func TestDeleteFamily_Execute(t *testing.T) {
	type args struct {
		auth *model.Auth
	}
	tests := []struct {
		name                     string
		mockUserRepoFunc         func(*repository.MockUserRepositoryInterface)
		mockFamilyRepoFunc       func(*repository.MockFamilyRepositoryInterface)
		mockMemberFamilyRepoFunc func(*repository.MockMemberFamilyRepositoryInterface)
		args                     args
		wantErr                  bool
		wantErrMsg               string
	}{
		{
			name: "success",
			mockUserRepoFunc: func(r *repository.MockUserRepositoryInterface) {
				r.EXPECT().GetUserWhereUsername(common.TestUserName).Return(&db.User{ID: common.TestUserID, Name: common.TestUserName}, nil)
			},
			mockFamilyRepoFunc: func(r *repository.MockFamilyRepositoryInterface) {
				r.EXPECT().DeleteFamily(common.TestFamilyID).Return(nil)
			},
			mockMemberFamilyRepoFunc: func(r *repository.MockMemberFamilyRepositoryInterface) {
				r.EXPECT().GetMemberFamilyWhereMemberID(common.TestUserID).Return(&db.MemberFamily{MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Role: common.TestRoleHead}, nil)
				r.EXPECT().DeleteMemberFamily(common.TestFamilyID).Return(nil)
			},
			args: args{
				auth: &model.Auth{UserName: common.TestUserName},
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
				familyRepo := repository.NewMockFamilyRepositoryInterface(ctrl)
				tt.mockFamilyRepoFunc(familyRepo)
				memberFamilyRepo := repository.NewMockMemberFamilyRepositoryInterface(ctrl)
				tt.mockMemberFamilyRepoFunc(memberFamilyRepo)

				// run target method
				err := NewDeleteFamily(userRepo, familyRepo, memberFamilyRepo).Execute(tt.args.auth)

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
