package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gold-kou/go-housework/app/server/service"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var successCreateUserReq = `
{
  "email": "test@example.com",
  "user_name": "test",
  "password": "123456"
}
`

var successCreateUserResp = `
{
  "user": {
    "user_id": 1,
    "user_name": "test"
  },
  "message": "new user created"
}
`

func TestCreateUser(t *testing.T) {
	type args struct {
		reqBody string
	}
	tests := []struct {
		name       string
		args       args
		dbMockFunc func(mock sqlmock.Sqlmock)
		mockFunc   func(*service.MockCreateUserServiceInterface)
		want       string
		wantStatus int
	}{
		/*
			{
				name: "success",
				args: args{reqBody: successCreateUserReq},
				dbMockFunc: func(mock sqlmock.Sqlmock) {
					mock.ExpectBegin()
				},
				mockFunc: func(c *service.MockCreateUserServiceInterface) {
					c.EXPECT().Execute().Return(&db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestHashedPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()}, nil)
				},
				want:       successCreateUserResp,
				wantStatus: 200,
			},

		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//var db *gorm.DB
			//common.SetDB(db)

			// db, _ := gorm.Open("postgres", "host=db port=5432 user=admin dbname=devdb password=admin! sslmode=disable")
			// common.SetDB(db)

			assert := assert.New(t)

			// mock service
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			createUserService := service.NewMockCreateUserServiceInterface(ctrl)
			tt.mockFunc(createUserService)

			// setting for execute
			req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(tt.args.reqBody))
			assert.NoError(err)
			resp := httptest.NewRecorder()

			// test target
			CreateUser(resp, req)

			// assert http code
			assert.Equal(tt.wantStatus, resp.Code)

			// assert http resp body
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(err)
			respBody := string(respBodyByte)
			assert.JSONEq(tt.want, respBody)
		})
	}
}
