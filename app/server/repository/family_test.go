package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gold-kou/go-housework/app/common"
	"github.com/gold-kou/go-housework/app/model/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestFamilyRepository_InsertFamily(t *testing.T) {
	type args struct {
		f *db.Family
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.Family
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "families" ("name","created_at","updated_at") VALUES (?,?,?)`)).
					WithArgs(common.TestFamilyName, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			args:    args{f: &db.Family{Name: common.TestFamilyName}},
			want:    &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &FamilyRepository{
					db: db,
				}

				// test target function
				err := r.InsertFamily(tt.args.f)

				// assert
				if tt.wantErr {
					assert.Error(err)
					assert.Equal(tt.wantErrMsg, err.Error())
				} else {
					assert.NoError(err)
					assert.Equal(tt.want, tt.args.f)
				}
			})
		})
	}
}

func TestFamilyRepository_UpdateFamily(t *testing.T) {
	type args struct {
		f *db.Family
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.Family
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				// omission
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "families"`)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			args:    args{f: &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName, CreatedAt: common.GetTestTime()}},
			want:    &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &FamilyRepository{
					db: db,
				}

				// test target function
				err := r.UpdateFamily(tt.args.f)

				// assert
				if tt.wantErr {
					assert.Error(err)
					assert.Equal(tt.wantErrMsg, err.Error())
				} else {
					assert.NoError(err)
					assert.Equal(tt.want, tt.args.f)
				}
			})
		})
	}
}

func TestFamilyRepository_DeleteFamily(t *testing.T) {
	type args struct {
		familyID uint64
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

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "families" WHERE (id = ?)`)).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			args:    args{familyID: common.TestFamilyID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &FamilyRepository{
					db: db,
				}

				// test target function
				err := r.DeleteFamily(tt.args.familyID)

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

func TestFamilyRepository_ShowFamily(t *testing.T) {
	type args struct {
		familyID uint64
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.Family
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "families" WHERE (id = ?)`)).
					WithArgs(common.TestFamilyID).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
						AddRow(common.TestFamilyID, common.TestFamilyName, common.GetGormTestTime(), common.GetGormTestTime()))
			},
			args:    args{familyID: common.TestFamilyID},
			want:    &db.Family{ID: common.TestFamilyID, Name: common.TestFamilyName, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &FamilyRepository{
					db: db,
				}

				// test target function
				got, err := r.ShowFamily(tt.args.familyID)

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
