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

var successCreateFamilyReqBody = `
{
  "family_name": "test-family"
}
`
var errCreateFamilyReqBodyEmptyFamilyName = `
{
}
`
var errCreateFamilyReqBodyWrongLengthFamilyName = `
{
  "family_name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`

var successCreateFamilyReq, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(successCreateFamilyReqBody))
var errCreateFamilyReqEmptyFamilyName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateFamilyReqBodyEmptyFamilyName))
var errCreateFamilyReqWrongLengthFamilyName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateFamilyReqBodyWrongLengthFamilyName))

var successCreateFamilyResp = &schemamodel.ResponseCreateFamily{
	Family: schemamodel.Family{
		FamilyId:   int64(common.TestFamilyID),
		FamilyName: common.TestFamilyName,
	},
}
var errCreateFamilyRespEmptyFamilyName = common.NewBadRequestError("family_name: cannot be blank.")
var errCreateFamilyRespWrongLengthFamilyName = common.NewBadRequestError("family_name: the length must be between 1 and 100.")

func TestCreateFamily(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockCreateFamilyServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successCreateFamilyReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateFamilyServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}, &schemamodel.RequestCreateFamily{FamilyName: common.TestFamilyName}).Return(&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, nil)
			},
			wantRespBody:   successCreateFamilyResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:       "error(validation error / family_name is empty)",
			testCaseID: 2,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateFamilyReqEmptyFamilyName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateFamilyServiceInterface) {},
			wantRespBody:    errCreateFamilyRespEmptyFamilyName,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / family_name is over length)",
			testCaseID: 3,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateFamilyReqWrongLengthFamilyName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateFamilyServiceInterface) {},
			wantRespBody:    errCreateFamilyRespWrongLengthFamilyName,
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
			createFamilyService := service.NewMockCreateFamilyServiceInterface(ctrl)
			tt.mockServiceFunc(createFamilyService)

			// test target
			h := CreateFamilyHandler{
				tok: tokenInterface,
				srv: createFamilyService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.CreateFamily(w, successCreateFamilyReq)
			case 2:
				resp, status, _ = h.CreateFamily(w, errCreateFamilyReqEmptyFamilyName)
			case 3:
				resp, status, _ = h.CreateFamily(w, errCreateFamilyReqWrongLengthFamilyName)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
