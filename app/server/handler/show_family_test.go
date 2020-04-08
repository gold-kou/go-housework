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

var successShowFamilyReq, _ = http.NewRequest(http.MethodPost, "", nil)

var successShowFamilyResp = &schemamodel.ResponseShowFamily{
	Family: schemamodel.Family{FamilyId: int64(common.TestFamilyID), FamilyName: common.TestFamilyName},
}

func TestShowFamily(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockShowFamilyServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successShowFamilyReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockShowFamilyServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}).Return(&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, nil)
			},
			wantRespBody:   successShowFamilyResp,
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
			ShowFamilyService := service.NewMockShowFamilyServiceInterface(ctrl)
			tt.mockServiceFunc(ShowFamilyService)

			// test target
			h := ShowFamilyHandler{
				tok: tokenInterface,
				srv: ShowFamilyService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.ShowFamily(w, successShowFamilyReq)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
