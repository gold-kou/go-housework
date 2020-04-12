package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/stretchr/testify/assert"
)

var successLoginReq, _ = http.NewRequest(http.MethodPost, "/login?user_name=test-user&password=123456", nil)
var errLoginReqEmptyUserName, _ = http.NewRequest(http.MethodPost, "/login?password=123456", nil)
var errLoginReqEmptyPassword, _ = http.NewRequest(http.MethodPost, "/login?user_name=test-user", nil)

var successLoginResp = &schemamodel.ResponseLogin{
	Token: common.TestToken,
}
var errLoginRespEmptyUserName = common.NewBadRequestError("cannot be blank")
var errLoginRespEmptyPassword = common.NewBadRequestError("cannot be blank")

func TestLogin(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockLoginServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:          "success",
			testCaseID:    1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {},
			mockServiceFunc: func(c *service.MockLoginServiceInterface) {
				c.EXPECT().Execute(common.TestUserName, common.TestPassword).Return(common.TestToken, nil)
			},
			wantRespBody:   successLoginResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:            "error(validation error / date is empty)",
			testCaseID:      2,
			mockTokenFunc:   func(m *middleware.MockTokenInterface) {},
			mockServiceFunc: func(c *service.MockLoginServiceInterface) {},
			wantRespBody:    errLoginRespEmptyUserName,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:            "error(validation error / date is wrong format)",
			testCaseID:      3,
			mockTokenFunc:   func(m *middleware.MockTokenInterface) {},
			mockServiceFunc: func(c *service.MockLoginServiceInterface) {},
			wantRespBody:    errLoginRespEmptyPassword,
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
			LoginService := service.NewMockLoginServiceInterface(ctrl)
			tt.mockServiceFunc(LoginService)

			// test target
			h := LoginHandler{
				tok: tokenInterface,
				srv: LoginService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.Login(w, successLoginReq)
			case 2:
				resp, status, _ = h.Login(w, errLoginReqEmptyUserName)
			case 3:
				resp, status, _ = h.Login(w, errLoginReqEmptyPassword)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
