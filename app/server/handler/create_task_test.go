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

var successCreateTaskReqBody = `
{
  "task_name": "test-task1",
  "member_name": "test-user",
  "status": "todo",
  "date": "2020-12-31"
}
`
var errCreateTaskReqBodyEmptyTaskName = `
{
  "member_name": "test-user",
  "status": "todo",
  "date": "2020-12-31"
}
`
var errCreateTaskReqBodyEmptyMemberName = `
{
  "task_name": "test-task1",
  "status": "todo",
  "date": "2020-12-31"
}
`
var errCreateTaskReqBodyEmptyStatus = `
{
  "task_name": "test-task1",
  "member_name": "test-user",
  "date": "2020-12-31"
}
`
var errCreateTaskReqBodyEmptyDate = `
{
  "task_name": "test-task1",
  "member_name": "test-user",
  "status": "todo"
}
`

// status is neither todo nor done
var errCreateTaskReqBodyStatusNotIn = `
{
  "task_name": "test-task1",
  "member_name": "test-user",
  "status": "a",
  "date": "2020-12-31"
}
`

// not correct date format
var errCreateTaskReqBodyDateWrongFormat = `
{
  "task_name": "test-task1",
  "member_name": "test-user",
  "status": "todo",
  "date": "a"
}
`

var successCreateTaskReq, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(successCreateTaskReqBody))
var errCreateTaskReqEmptyTaskName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateTaskReqBodyEmptyTaskName))
var errCreateTaskReqEmptyMemberName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateTaskReqBodyEmptyMemberName))
var errCreateTaskReqEmptyStatus, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateTaskReqBodyEmptyStatus))
var errCreateTaskReqEmptyDate, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateTaskReqBodyEmptyDate))
var errCreateTaskReqStatusNotIn, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateTaskReqBodyStatusNotIn))
var errCreateTaskReqDateWrongFormat, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errCreateTaskReqBodyDateWrongFormat))

var successCreateTaskResp = &schemamodel.ResponseCreateTask{
	Family: schemamodel.Family{FamilyId: int64(common.TestFamilyID), FamilyName: common.TestFamilyName},
	Task:   schemamodel.Task{TaskId: int64(common.TestTaskID1), TaskName: common.TestTaskName1, MemberName: common.TestUserName, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate},
}
var errCreateTaskRespEmptyTaskName = common.NewBadRequestError("task_name: cannot be blank.")
var errCreateTaskRespEmptyMemberName = common.NewBadRequestError("member_name: cannot be blank.")
var errCreateTaskRespStatus = common.NewBadRequestError("status: cannot be blank.")
var errCreateTaskRespEmptyDate = common.NewBadRequestError("date: cannot be blank.")
var errCreateTaskRespStatusNotIn = common.NewBadRequestError("status: must be a valid value.")
var errCreateTaskRespDateWrongFormat = common.NewBadRequestError("date: must be a valid date.")

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockCreateTaskServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successCreateTaskReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName},
					&schemamodel.RequestCreateTask{TaskName: common.TestTaskName1, MemberName: common.TestUserName, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate}).
					Return(&db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword},
						&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName},
						&db.Task{ID: common.TestTaskID1, Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate}, nil)
			},
			wantRespBody:   successCreateTaskResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:       "error(validation error / task_name is empty)",
			testCaseID: 2,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateTaskReqEmptyTaskName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {},
			wantRespBody:    errCreateTaskRespEmptyTaskName,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / member_name is empty)",
			testCaseID: 3,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateTaskReqEmptyMemberName).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {},
			wantRespBody:    errCreateTaskRespEmptyMemberName,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / status is empty)",
			testCaseID: 4,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateTaskReqEmptyStatus).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {},
			wantRespBody:    errCreateTaskRespStatus,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / date is empty)",
			testCaseID: 5,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateTaskReqEmptyDate).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {},
			wantRespBody:    errCreateTaskRespEmptyDate,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / status is not in todo or done)",
			testCaseID: 6,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateTaskReqStatusNotIn).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {},
			wantRespBody:    errCreateTaskRespStatusNotIn,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:       "error(validation error / date is wrong format)",
			testCaseID: 7,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(errCreateTaskReqDateWrongFormat).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockCreateTaskServiceInterface) {},
			wantRespBody:    errCreateTaskRespDateWrongFormat,
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
			CreateTaskService := service.NewMockCreateTaskServiceInterface(ctrl)
			tt.mockServiceFunc(CreateTaskService)

			// test target
			h := CreateTaskHandler{
				tok: tokenInterface,
				srv: CreateTaskService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.CreateTask(w, successCreateTaskReq)
			case 2:
				resp, status, _ = h.CreateTask(w, errCreateTaskReqEmptyTaskName)
			case 3:
				resp, status, _ = h.CreateTask(w, errCreateTaskReqEmptyMemberName)
			case 4:
				resp, status, _ = h.CreateTask(w, errCreateTaskReqEmptyStatus)
			case 5:
				resp, status, _ = h.CreateTask(w, errCreateTaskReqEmptyDate)
			case 6:
				resp, status, _ = h.CreateTask(w, errCreateTaskReqStatusNotIn)
			case 7:
				resp, status, _ = h.CreateTask(w, errCreateTaskReqDateWrongFormat)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
