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

var successUpdateFamilyReqBody = `
{
  "family_name": "test-family"
}
`
var errUpdateFamilyReqBodyEmptyFamilyName = `
{
}
`
var errUpdateFamilyReqBodyWrongLengthFamilyName = `
{
  "family_name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`

var successUpdateFamilyReq, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(successUpdateFamilyReqBody))
var errUpdateFamilyReqEmptyFamilyName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateFamilyReqBodyEmptyFamilyName))
var errUpdateFamilyReqWrongLengthFamilyName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateFamilyReqBodyWrongLengthFamilyName))

var successUpdateFamilyResp = &schemamodel.ResponseUpdateFamily{
	Family: schemamodel.Family{
		FamilyId:   int64(common.TestFamilyID),
		FamilyName: common.TestFamilyName,
	},
}
var errUpdateFamilyRespEmptyFamilyName = common.NewBadRequestError("family_name: cannot be blank.")
var errUpdateFamilyRespWrongLengthFamilyName = common.NewBadRequestError("family_name: the length must be between 1 and 100.")

func TestUpdateFamily(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockUpdateFamilyServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successUpdateFamilyReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockUpdateFamilyServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}, &schemamodel.RequestUpdateFamily{FamilyName: common.TestFamilyName}).Return(&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, nil)
			},
			wantRespBody:   successUpdateFamilyResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:       "error(validation error / family_name is empty)",
			testCaseID: 2,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errUpdateFamilyReqEmptyFamilyName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockUpdateFamilyServiceInterface) {},
			wantRespBody:    errUpdateFamilyRespEmptyFamilyName,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / family_name is over length)",
			testCaseID: 3,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errUpdateFamilyReqWrongLengthFamilyName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockUpdateFamilyServiceInterface) {},
			wantRespBody:    errUpdateFamilyRespWrongLengthFamilyName,
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
			UpdateFamilyService := service.NewMockUpdateFamilyServiceInterface(ctrl)
			tt.mockServiceFunc(UpdateFamilyService)

			// test target
			h := UpdateFamilyHandler{
				tok: tokenInterface,
				srv: UpdateFamilyService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.UpdateFamily(w, successUpdateFamilyReq)
			case 2:
				resp, status, _ = h.UpdateFamily(w, errUpdateFamilyReqEmptyFamilyName)
			case 3:
				resp, status, _ = h.UpdateFamily(w, errUpdateFamilyReqWrongLengthFamilyName)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
