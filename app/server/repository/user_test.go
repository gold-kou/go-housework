package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_InsertUser(t *testing.T) {
	type args struct {
		u *db.User
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","password","created_at","updated_at") VALUES (?,?,?,?,?)`)).
					WithArgs(common.TestUserName, common.TestEmail, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			args:    args{u: &db.User{Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail}},
			want:    &db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
		{
			name: "fail(duplicate user name)",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","password","created_at","updated_at") VALUES (?,?,?,?,?)`)).
					WithArgs(common.TestUserName, common.TestEmail, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"users_name_key\""))

				mock.ExpectRollback()
			},
			args:       args{u: &db.User{Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail}},
			want:       &db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr:    true,
			wantErrMsg: "pq: duplicate key value violates unique constraint \"users_name_key\"",
		},
		{
			name: "fail(duplicate email)",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("name","email","password","created_at","updated_at") VALUES (?,?,?,?,?)`)).
					WithArgs(common.TestUserName, common.TestEmail, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"users_email_key\""))

				mock.ExpectRollback()
			},
			args:       args{u: &db.User{Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail}},
			want:       &db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr:    true,
			wantErrMsg: "pq: duplicate key value violates unique constraint \"users_email_key\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set test time
			common.SetTestTime()
			defer common.ResetTime()
			common.SetGormTestTime()
			defer common.ResetGormTime()

			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &UserRepository{
					db: db,
				}

				// test target function
				err := r.InsertUser(tt.args.u)

				// assert
				if tt.wantErr {
					assert.Error(err)
					assert.Equal(tt.wantErrMsg, err.Error())
				} else {
					assert.NoError(err)
					assert.Equal(tt.want, tt.args.u)
				}
			})
		})
	}
}

func TestUserRepository_GetUserWhereUsername(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (name = ?)`)).
					WithArgs(common.TestUserName).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
						AddRow(common.TestUserID, common.TestUserName, common.TestEmail, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()))
			},
			args:    args{userName: common.TestUserName},
			want:    &db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fix test time
			common.SetTestTime()
			defer common.ResetTime()
			common.SetGormTestTime()
			defer common.ResetGormTime()

			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &UserRepository{
					db: db,
				}

				// test target function
				got, err := r.GetUserWhereUsername(tt.args.userName)

				// assert
				if tt.wantErr {
					assert.Error(err)
					assert.Equal(tt.wantErrMsg, err.Error())
				} else {
					assert.NoError(err)
					assert.Equal(tt.want, got)
				}
			})
		})
	}
}

func TestUserRepository_DeleteUserWhereUsername(t *testing.T) {
	type args struct {
		userName string
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE (name = ?)`)).
					WithArgs(common.TestUserName).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			args:    args{userName: common.TestUserName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set test time
			common.SetTestTime()
			defer common.ResetTime()
			common.SetGormTestTime()
			defer common.ResetGormTime()

			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &UserRepository{
					db: db,
				}

				// test target function
				err := r.DeleteUserWhereUsername(tt.args.userName)

				// assert
				if tt.wantErr {
					assert.Error(err)
					assert.Equal(tt.wantErrMsg, err.Error())
				} else {
					assert.NoError(err)
				}
			})
		})
	}
}
