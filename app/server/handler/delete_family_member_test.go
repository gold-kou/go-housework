package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var successDeleteFamilyMemberReq, _ = http.NewRequest(http.MethodPost, "", nil)

var successDeleteFamilyMemberResp = &schemamodel.ResponseDeleteFamilyMember{
	Message: "delete complete",
}

func TestDeleteFamilyMember(t *testing.T) {
	// set path parameter
	successDeleteFamilyMemberReq = mux.SetURLVars(successDeleteFamilyMemberReq, map[string]string{"member_id": strconv.Itoa(int(common.TestUserID))})

	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockDeleteFamilyMemberServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successDeleteFamilyMemberReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockDeleteFamilyMemberServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}, common.TestUserID).Return(nil)
			},
			wantRespBody:   successDeleteFamilyMemberResp,
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
			DeleteFamilyMemberService := service.NewMockDeleteFamilyMemberServiceInterface(ctrl)
			tt.mockServiceFunc(DeleteFamilyMemberService)

			// test target
			h := DeleteFamilyMemberHandler{
				tok: tokenInterface,
				srv: DeleteFamilyMemberService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.DeleteFamilyMember(w, successDeleteFamilyMemberReq)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
