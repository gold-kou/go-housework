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

var successUpdateTaskReqBody = `
{
  "task": {
    "task_id": 1,
    "task_name": "test-task1",
    "member_name": "test-user",
    "status": "todo",
    "date": "2020-12-31"
  }
}
`
var errUpdateTaskReqBodyEmptyTaskID = `
{
  "task": {
    "task_name": "test-task1",
    "member_name": "test-user",
    "status": "todo",
    "date": "2020-12-31"
  }
}
`
var errUpdateTaskReqBodyEmptyTaskName = `
{
  "task": {
    "task_id": 1,
    "member_name": "test-user",
    "status": "todo",
    "date": "2020-12-31"
  }
}
`
var errUpdateTaskReqBodyEmptyMemberName = `
{
  "task": {
    "task_id": 1,
    "task_name": "test-task1",
    "status": "todo",
    "date": "2020-12-31"
  }
}
`
var errUpdateTaskReqBodyEmptyStatus = `
{
  "task": {
    "task_id": 1,
    "task_name": "test-task1",
    "member_name": "test-user",
    "date": "2020-12-31"
  }
}
`
var errUpdateTaskReqBodyEmptyDate = `
{
  "task": {
    "task_id": 1,
    "task_name": "test-task1",
    "member_name": "test-user",
    "status": "todo"
  }
}
`

// status is neither todo nor done
var errUpdateTaskReqBodyStatusNotIn = `
{
  "task": {
    "task_id": 1,
    "task_name": "test-task1",
    "member_name": "test-user",
    "status": "a",
    "date": "2020-12-31"
  }
}
`

// not correct date format
var errUpdateTaskReqBodyDateWrongFormat = `
{
  "task": {
    "task_id": 1,
    "task_name": "test-task1",
    "member_name": "test-user",
    "status": "todo",
    "date": "a"
  }
}
`

var successUpdateTaskReq, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(successUpdateTaskReqBody))
var errUpdateTaskReqEmptyTaskID, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyEmptyTaskID))
var errUpdateTaskReqEmptyTaskName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyEmptyTaskName))
var errUpdateTaskReqEmptyMemberName, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyEmptyMemberName))
var errUpdateTaskReqEmptyStatus, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyEmptyStatus))
var errUpdateTaskReqEmptyDate, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyEmptyDate))
var errUpdateTaskReqStatusNotIn, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyStatusNotIn))
var errUpdateTaskReqDateWrongFormat, _ = http.NewRequest(http.MethodPost, "", strings.NewReader(errUpdateTaskReqBodyDateWrongFormat))

var successUpdateTaskResp = &schemamodel.ResponseUpdateTask{
	Family: schemamodel.Family{FamilyId: int64(common.TestFamilyID), FamilyName: common.TestFamilyName},
	Task:   schemamodel.Task{TaskId: int64(common.TestTaskID1), TaskName: common.TestTaskName1, MemberName: common.TestUserName, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate},
}
var errUpdateTaskRespEmptyTaskID = common.NewBadRequestError("task_id: cannot be blank.")
var errUpdateTaskRespEmptyTaskName = common.NewBadRequestError("task_name: cannot be blank.")
var errUpdateTaskRespEmptyMemberName = common.NewBadRequestError("member_name: cannot be blank.")
var errUpdateTaskRespStatus = common.NewBadRequestError("status: cannot be blank.")
var errUpdateTaskRespEmptyDate = common.NewBadRequestError("date: cannot be blank.")
var errUpdateTaskRespStatusNotIn = common.NewBadRequestError("status: must be a valid value.")
var errUpdateTaskRespDateWrongFormat = common.NewBadRequestError("date: must be a valid date.")

func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name            string
		testCaseID      int
		reqType         string
		mockTokenFunc   func(*middleware.MockTokenInterface)
		mockServiceFunc func(*service.MockUpdateTaskServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:       "success",
			testCaseID: 1,
			mockTokenFunc: func(m *middleware.MockTokenInterface) {
				m.EXPECT().VerifyHeaderToken(successUpdateTaskReq).Return(&model.Auth{UserName: common.TestUserName}, nil)
			},
			mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {
				c.EXPECT().Execute(&model.Auth{UserName: common.TestUserName},
					&schemamodel.RequestUpdateTask{Task: schemamodel.Task{TaskId: int64(common.TestTaskID1), TaskName: common.TestTaskName1, MemberName: common.TestUserName, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate}}).
					Return(&db.Task{ID: common.TestTaskID1, Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate},
						&db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName},
						&db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword}, nil)
			},
			wantRespBody:   successUpdateTaskResp,
			wantRespStatus: http.StatusOK,
		},
		/*
			{
				name:       "error(validation error / task_id is empty)",
				testCaseID: 2,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqEmptyTaskID).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespEmptyTaskID,
				wantRespStatus:  http.StatusBadRequest,
			},
			{
				name:       "error(validation error / task_name is empty)",
				testCaseID: 3,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqEmptyTaskName).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespEmptyTaskName,
				wantRespStatus:  http.StatusBadRequest,
			},
			{
				name:       "error(validation error / member_name is empty)",
				testCaseID: 4,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqEmptyMemberName).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespEmptyMemberName,
				wantRespStatus:  http.StatusBadRequest,
			},
			{
				name:       "error(validation error / status is empty)",
				testCaseID: 5,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqEmptyStatus).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespStatus,
				wantRespStatus:  http.StatusBadRequest,
			},
			{
				name:       "error(validation error / date is empty)",
				testCaseID: 6,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqEmptyDate).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespEmptyDate,
				wantRespStatus:  http.StatusBadRequest,
			},
			{
				name:       "error(validation error / status is not in todo or done)",
				testCaseID: 7,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqStatusNotIn).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespStatusNotIn,
				wantRespStatus:  http.StatusBadRequest,
			},
			{
				name:       "error(validation error / date is wrong format)",
				testCaseID: 8,
				mockTokenFunc: func(m *middleware.MockTokenInterface) {
					m.EXPECT().VerifyHeaderToken(errUpdateTaskReqDateWrongFormat).Return(&model.Auth{UserName: common.TestUserName}, nil)
				},
				mockServiceFunc: func(c *service.MockUpdateTaskServiceInterface) {},
				wantRespBody:    errUpdateTaskRespDateWrongFormat,
				wantRespStatus:  http.StatusBadRequest,
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			// mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tokenInterface := middleware.NewMockTokenInterface(ctrl)
			tt.mockTokenFunc(tokenInterface)
			UpdateTaskService := service.NewMockUpdateTaskServiceInterface(ctrl)
			tt.mockServiceFunc(UpdateTaskService)

			// test target
			h := UpdateTaskHandler{
				tok: tokenInterface,
				srv: UpdateTaskService,
			}
			w := httptest.NewRecorder()
			var resp interface{}
			var status int
			// change request case for test patterns
			switch tt.testCaseID {
			case 1:
				resp, status, _ = h.UpdateTask(w, successUpdateTaskReq)
			case 2:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqEmptyTaskID)
			case 3:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqEmptyTaskName)
			case 4:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqEmptyMemberName)
			case 5:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqEmptyStatus)
			case 6:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqEmptyDate)
			case 7:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqStatusNotIn)
			case 8:
				resp, status, _ = h.UpdateTask(w, errUpdateTaskReqDateWrongFormat)
			}

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
