package handler

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/middleware"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/stretchr/testify/assert"
)

var successDeleteTaskReq, _ = http.NewRequest(http.MethodPost, "", nil)

var successDeleteTaskResp = &schemamodel.ResponseDeleteTask{
	Message: "delete complete",
}

func TestDeleteTask(t *testing.T) {
	// set path parameter
	successDeleteTaskReq = mux.SetURLVars(successDeleteTaskReq, map[string]string{"task_id": strconv.Itoa(int(common.TestTaskID1))})

	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockDeleteTaskServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successDeleteTaskReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockDeleteTaskServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}, common.TestTaskID1).Return(nil)
			},
			wantRespBody:   successDeleteTaskResp,
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
			DeleteTaskService := service.NewMockDeleteTaskServiceInterface(ctrl)
			tt.mockServiceFunc(DeleteTaskService)

			// test target
			h := DeleteTaskHandler{
				tok: tokenInterface,
				srv: DeleteTaskService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.DeleteTask(w, successDeleteTaskReq)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
