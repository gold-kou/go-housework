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

var successListTasksReq, _ = http.NewRequest(http.MethodPost, "/tasks?date=2020-12-31", nil)
var errListTasksReqEmptyDate, _ = http.NewRequest(http.MethodPost, "/tasks", nil)
var errListTasksReqDateWrongFormat, _ = http.NewRequest(http.MethodPost, "/tasks?date=a", nil)

var successListTasksResp = &schemamodel.ResponseListTasks{
	Family: schemamodel.Family{FamilyId: int64(common.TestFamilyID), FamilyName: common.TestFamilyName},
	Tasks:  []schemamodel.Task{{TaskId: int64(common.TestTaskID1), TaskName: common.TestTaskName1, MemberName: common.TestUserName, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate}},
}
var errListTasksRespEmptyDate = common.NewBadRequestError("cannot be blank")
var errListTasksRespDateWrongFormat = common.NewBadRequestError("must be a valid date")

func TestListTasks(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockListTasksServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successListTasksReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockListTasksServiceInterface) {
				var ts []*db.Task
				ts = append(ts, &db.Task{ID: common.TestTaskID1, Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate})
				var us []*db.User
				us = append(us, &db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword})
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName}, common.TestTaskDate).Return(ts, &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName}, us, nil)
			},
			wantRespBody:   successListTasksResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:       "error(validation error / date is empty)",
			testCaseID: 2,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errListTasksReqEmptyDate).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockListTasksServiceInterface) {},
			wantRespBody:    errListTasksRespEmptyDate,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / date is wrong format)",
			testCaseID: 3,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errListTasksReqDateWrongFormat).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockListTasksServiceInterface) {},
			wantRespBody:    errListTasksRespDateWrongFormat,
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
			ListTasksService := service.NewMockListTasksServiceInterface(ctrl)
			tt.mockServiceFunc(ListTasksService)

			// test target
			h := ListTasksHandler{
				tok: tokenInterface,
				srv: ListTasksService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.ListTasks(w, successListTasksReq)
			case 2:
				resp, status, _ = h.ListTasks(w, errListTasksReqEmptyDate)
			case 3:
				resp, status, _ = h.ListTasks(w, errListTasksReqDateWrongFormat)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
