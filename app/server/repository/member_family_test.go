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

func TestMemberFamilyRepository_InsertMemberFamily(t *testing.T) {
	type args struct {
		mf *db.MemberFamily
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success(role is head)",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "members_families" ("member_id","family_id","role","created_at","updated_at") VALUES (?,?,?,?,?)`)).
					WithArgs(common.TestUserID, common.TestFamilyID, db.FamilyRoleHead, common.GetGormTestTime(), common.GetGormTestTime()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			args:    args{mf: &db.MemberFamily{MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Role: db.FamilyRoleHead}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &MemberFamilyRepository{
					db: db,
				}

				// test target function
				err := r.InsertMemberFamily(tt.args.mf)

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

func TestMemberFamilyRepository_DeleteMemberFamily(t *testing.T) {
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

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "members_families" WHERE (family_id = ?)`)).
					WithArgs(common.TestFamilyID).
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
				r := &MemberFamilyRepository{
					db: db,
				}

				// test target function
				err := r.DeleteMemberFamily(tt.args.familyID)

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

func TestMemberFamilyRepository_DeleteMemberFromFamily(t *testing.T) {
	type args struct {
		memberID uint64
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

				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "members_families" WHERE (member_id = ?)`)).
					WithArgs(common.TestFamilyID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			args:    args{memberID: common.TestUserID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &MemberFamilyRepository{
					db: db,
				}

				// test target function
				err := r.DeleteMemberFromFamily(tt.args.memberID)

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

func TestFamilyRepository_GetMemberFamilyWhereMemberID(t *testing.T) {
	type args struct {
		memberID uint64
	}
	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       *db.MemberFamily
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "members_families" WHERE (member_id = ?)`)).
					WithArgs(common.TestUserID).
					WillReturnRows(sqlmock.NewRows([]string{"member_id", "family_id", "role", "created_at", "updated_at"}).
						AddRow(common.TestUserID, common.TestFamilyID, db.FamilyRoleHead, common.GetGormTestTime(), common.GetGormTestTime()))
			},
			args:    args{memberID: common.TestUserID},
			want:    &db.MemberFamily{MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Role: db.FamilyRoleHead, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &MemberFamilyRepository{
					db: db,
				}

				// test target function
				got, err := r.GetMemberFamilyWhereMemberID(tt.args.memberID)

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

func TestFamilyRepository_ListMemberFamiliesWhereFamilyID(t *testing.T) {
	type args struct {
		familyID uint64
	}
	want := make([]*db.MemberFamily, 0)
	want = append(want, &db.MemberFamily{MemberID: common.TestUserID, FamilyID: common.TestFamilyID, Role: db.FamilyRoleHead, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()})
	want = append(want, &db.MemberFamily{MemberID: common.TestUserID2, FamilyID: common.TestFamilyID, Role: db.FamilyRoleMember, CreatedAt: common.GetTestTime(), UpdatedAt: common.GetTestTime()})

	tests := []struct {
		name       string
		dbMockFunc func(mock sqlmock.Sqlmock)
		args       args
		want       []*db.MemberFamily
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "success",
			dbMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "members_families" WHERE (family_id = ?)`)).
					WithArgs(common.TestUserID).
					WillReturnRows(sqlmock.NewRows([]string{"member_id", "family_id", "role", "created_at", "updated_at"}).
						AddRow(common.TestUserID, common.TestFamilyID, db.FamilyRoleHead, common.GetGormTestTime(), common.GetGormTestTime()).
						AddRow(common.TestUserID2, common.TestFamilyID, db.FamilyRoleMember, common.GetGormTestTime(), common.GetGormTestTime()))
			},
			args:    args{familyID: common.TestFamilyID},
			want:    want,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			common.MockDB(t, func(db *gorm.DB, mock sqlmock.Sqlmock) {
				tt.dbMockFunc(mock)
				r := &MemberFamilyRepository{
					db: db,
				}

				// test target function
				got, err := r.ListMemberFamiliesWhereFamilyID(tt.args.familyID)

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
