package service

// TODO test

/*
import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/repository"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestCreateFamily_Execute(t *testing.T) {
	type args struct {
		auth *middleware.Auth
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
				reqCreateFamily := &schemamodel.RequestCreateFamily{FamilyName: common.TestFamilyName}
				_, err := NewCreateFamily(db, reqCreateFamily, familyRepo, userRepo, *memberFamilyRepo).Execute(tt.args.auth)

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
*/
