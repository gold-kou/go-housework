package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_InsertTask(t *testing.T) {
	type args struct {
		t *db.Task
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tasks" ("name","member_id","family_id","status","date","created_at","updated_at") VALUES (?,?,?,?,?,?,?)`)).
					WithArgs(common.TestTaskName1, common.TestUserID, common.TestFamilyID, common.TestTaskStatusTodo, common.TestTaskDate, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			args:    args{t: &db.Task{Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()}},
			wantErr: false,
		},
		{
			name: "fail(duplicate task)",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "tasks" ("name","member_id","family_id","status","date","created_at","updated_at") VALUES (?,?,?,?,?,?,?)`)).
					WithArgs(common.TestTaskName1, common.TestUserID, common.TestFamilyID, common.TestTaskStatusTodo, common.TestTaskDate, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnError(errors.New("pq: duplicate key value violates unique constraint \"uk_tasks\""))

				mock.ExpectRollback()
			},
			args:       args{t: &db.Task{Name: common.TestTaskName1, MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Status: common.TestTaskStatusTodo, Date: common.TestTaskDate, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()}},
			wantErr:    true,
			wantErrMsg: "pq: duplicate key value violates unique constraint \"uk_tasks\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &TaskRepository{
					db: db,
				}

				// test target function
				err := r.InsertTask(tt.args.t)

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
