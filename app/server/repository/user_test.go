package repository

import (
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// TestMain prepare common test setting
func TestMain(m *testing.M) {
	os.Exit(common.RunTestMain(m))
}

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

func TestUserRepository_GetUserWhereUserID(t *testing.T) {
	type args struct {
		userID uint64
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
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (id = ?)`)).
					WithArgs(common.TestUserID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
						AddRow(common.TestUserID, common.TestUserName, common.TestEmail, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()))
			},
			args:    args{userID: common.TestUserID},
			want:    &db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &UserRepository{
					db: db,
				}

				// test target function
				got, err := r.GetUserWhereUserID(tt.args.userID)

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

func TestUserRepository_GetUsersWhereUserIDs(t *testing.T) {
	type args struct {
		userIDs []uint64
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       []*db.User
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (id = ?)`)).
					WithArgs(common.TestUserID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
						AddRow(common.TestUserID, common.TestUserName, common.TestEmail, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (id = ?)`)).
					WithArgs(common.TestUserID2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
						AddRow(common.TestUserID2, common.TestUserName2, common.TestEmail2, common.TestPassword, common.GetGormTestTime(), common.GetGormTestTime()))
			},
			args: args{userIDs: []uint64{common.TestUserID, common.TestUserID2}},
			want: []*db.User{
				&db.User{ID: common.TestUserID, Name: common.TestUserName, Password: common.TestPassword, Email: common.TestEmail, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
				&db.User{ID: common.TestUserID2, Name: common.TestUserName2, Password: common.TestPassword, Email: common.TestEmail2, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &UserRepository{
					db: db,
				}

				// test target function
				got, err := r.GetUsersWhereUserIDs(tt.args.userIDs)

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
