package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gold-kou/go-housework/app/model/db"

	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/stretchr/testify/assert"
)

var successListFamilyMembersReq, _ = http.NewRequest(http.MethodPost, "", nil)

var successListFamilyMembersResp = &schemamodel.ResponseListFamilyMembers{
	Family:  schemamodel.Family{FamilyId: int64(common.TestFamilyID), FamilyName: common.TestFamilyName},
	Members: []schemamodel.Member{{MemberId: int64(common.TestUserID), MemberName: common.TestUserName}},
}

func TestListFamilyMembers(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockListFamilyMembersServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successListFamilyMembersReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockListFamilyMembersServiceInterface) {
				var us []*db.User
				us = append(us, &db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword})
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}).Return(&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, us, nil)
			},
			wantRespBody:   successListFamilyMembersResp,
			wantRespStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			// mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tokenInterface := middleware.NewMockTokenInterface(ctrl)
			tt.mockTokenFunc(tokenInterface)
			ListFamilyMembersService := service.NewMockListFamilyMembersServiceInterface(ctrl)
			tt.mockServiceFunc(ListFamilyMembersService)

			// test target
			h := ListFamilyMembersHandler{
				tok: tokenInterface,
				srv: ListFamilyMembersService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.ListFamilyMembers(w, successListFamilyMembersReq)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
