package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/gold-kou/go-housework/app/model/schemamodel"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/stretchr/testify/assert"
)

var successCreateUserReq = `
{
  "email": "test@example.com",
  "user_name": "test-user",
  "password": "123456"
}
`
var errCreateUserReqEmptyEmail = `
{
  "user_name": "test-user",
  "password": "123456"
}
`
var errCreateUserReqEmptyUserName = `
{
  "email": "test@example.com",
  "password": "123456"
}
`
var errCreateUserReqWrongLengthUserName = `
{
  "email": "test@example.com",
  "user_name": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
  "password": "123456"
}
`
var errCreateUserReqEmptyPassword = `
{
  "email": "test@example.com",
  "user_name": "test-user"
}
`
var errCreateUserReqWrongLengthPassword = `
{
  "email": "test@example.com",
  "user_name": "test-user",
  "password": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}
`

var successCreateUserResp = &schemamodel.ResponseCreateUser{
	User: schemamodel.User{
		UserId:   int64(common.TestUserID),
		UserName: common.TestUserName,
	},
}
var errCreateUserRespEmptyEmail = common.NewBadRequestError("email: cannot be blank.")
var errCreateUserRespEmptyUserName = common.NewBadRequestError("user_name: cannot be blank.")
var errCreateUserRespWrongLengthUserName = common.NewBadRequestError("user_name: the length must be between 1 and 100.")
var errCreateUserRespEmptyPassword = common.NewBadRequestError("password: cannot be blank.")
var errCreateUserRespWrongLengthPassword = common.NewBadRequestError("password: the length must be between 6 and 20.")

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name            string
		reqBody         string
		mockServiceFunc func(*service.MockCreateUserServiceInterface)
		wantRespBody    interface{}
		wantRespStatus  int
	}{
		{
			name:    "success",
			reqBody: successCreateUserReq,
			mockServiceFunc: func(c *service.MockCreateUserServiceInterface) {
				c.EXPECT().Execute(&schemamodel.RequestCreateUser{Email: common.TestEmail, UserName: common.TestUserName, Password: common.TestPassword}).Return(&db.User{ID: common.TestUserID, Name: common.TestUserName, Email: common.TestEmail, Password: common.TestHashedPassword}, nil)
			},
			wantRespBody:   successCreateUserResp,
			wantRespStatus: http.StatusOK,
		},
		{
			name:            "failed(validation error / email is empty)",
			reqBody:         errCreateUserReqEmptyEmail,
			mockServiceFunc: func(c *service.MockCreateUserServiceInterface) {},
			wantRespBody:    errCreateUserRespEmptyEmail,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:            "failed(validation error / user_name is empty)",
			reqBody:         errCreateUserReqEmptyEmail,
			mockServiceFunc: func(c *service.MockCreateUserServiceInterface) {},
			wantRespBody:    errCreateUserRespEmptyEmail,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:            "failed(validation error / user_name is over length)",
			reqBody:         errCreateUserReqWrongLengthUserName,
			mockServiceFunc: func(c *service.MockCreateUserServiceInterface) {},
			wantRespBody:    errCreateUserRespWrongLengthUserName,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:            "failed(validation error / password is empty)",
			reqBody:         errCreateUserReqEmptyPassword,
			mockServiceFunc: func(c *service.MockCreateUserServiceInterface) {},
			wantRespBody:    errCreateUserRespEmptyPassword,
			wantRespStatus:  http.StatusBadRequest,
		},
		{
			name:            "failed(validation error / password is over length)",
			reqBody:         errCreateUserReqWrongLengthPassword,
			mockServiceFunc: func(c *service.MockCreateUserServiceInterface) {},
			wantRespBody:    errCreateUserRespWrongLengthPassword,
			wantRespStatus:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			// mock
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			createUserService := service.NewMockCreateUserServiceInterface(ctrl)
			tt.mockServiceFunc(createUserService)

			h := CreateUserHandler{
				srv: createUserService,
			}

			// test target
			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(tt.reqBody))
			assert.NoError(err)
			resp, status, _ := h.CreateUser(w, r)

			// assert http response body
			assert.Equal(tt.wantRespBody, resp)

			// assert http status code
			assert.Equal(tt.wantRespStatus, status)
		})
	}
}
