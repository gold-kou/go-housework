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

func TestDeleteFamilyMember_Execute(t *testing.T) {
	type args struct {
		auth     *model.Auth
		memberID uint64
	}
	tests := []struct {
		name                     string
		mockMemberFamilyRepoFunc func(*repository.MockMemberFamilyRepositoryInterface)
		args                     args
		wantErr                  bool
		wantErrMsg               string
	}{
		{
			name: "success",
			mockMemberFamilyRepoFunc: func(r *repository.MockMemberFamilyRepositoryInterface) {
				r.EXPECT().DeleteMemberFromFamily(common.TestUserID).Return(nil)
			},
			args: args{
				auth:     &model.Auth{UserName: common.TestUserName},
				memberID: common.TestUserID,
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
				memberFamilyRepo := repository.NewMockMemberFamilyRepositoryInterface(ctrl)
				tt.mockMemberFamilyRepoFunc(memberFamilyRepo)

				// run target method
				err := NewDeleteFamilyMember(memberFamilyRepo).Execute(tt.args.auth, tt.args.memberID)

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
