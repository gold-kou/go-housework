package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/stretchr/testify/assert"
)

var successRegisterFamilyMemberReqBody = `
{
  "member_name": "test-user"
}
`
var errRegisterFamilyMemberReqBodyEmptyMemberName = `
{
}
`

var successRegisterFamilyMemberReq, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(successRegisterFamilyMemberReqBody))
var errRegisterFamilyMemberReqEmptyMemberName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errRegisterFamilyMemberReqBodyEmptyMemberName))

var successRegisterFamilyMemberResp = &schemamodel.ResponseRegisterFamilyMember{
	Family: schemamodel.Family{
		FamilyId:   int64(common.TestFamilyID),
		FamilyName: common.TestFamilyName,
	},
	Member: schemamodel.Member{
		MemberId:   int64(common.TestUserID),
		MemberName: common.TestUserName,
	},
}
var errRegisterFamilyMemberRespEmptyMemberName = common.NewBadRequestError("member_name: cannot be blank.")

func TestRegisterFamilyMember(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockRegisterFamilyMemberServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successRegisterFamilyMemberReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockRegisterFamilyMemberServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}, &schemamodel.RequestRegisterFamilyMember{MemberName: common.TestUserName}).
					Return(&db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword}, &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, nil)
			},
			wantRespBody:   successRegisterFamilyMemberResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:       "error(validation error / member_name is empty)",
			testCaseID: 2,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errRegisterFamilyMemberReqEmptyMemberName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockRegisterFamilyMemberServiceInterface) {},
			wantRespBody:    errRegisterFamilyMemberRespEmptyMemberName,
			wantRespStatus:  http.StatusBadRequest,
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
			RegisterFamilyMemberService := service.NewMockRegisterFamilyMemberServiceInterface(ctrl)
			tt.mockServiceFunc(RegisterFamilyMemberService)

			// test target
			h := RegisterFamilyMemberHandler{
				tok: tokenInterface,
				srv: RegisterFamilyMemberService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.RegisterFamilyMember(w, successRegisterFamilyMemberReq)
			case 2:
				resp, status, _ = h.RegisterFamilyMember(w, errRegisterFamilyMemberReqEmptyMemberName)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
