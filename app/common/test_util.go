package common

import (
	"database/sql/driver"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// TestUserID TestUserID
	TestUserID = uint64(1)
	// TestUserID2 TestUserID2
	TestUserID2 = uint64(2)
	// TestUserName TestUserName
	TestUserName = "test-user"
	// TestUserName2 TestUserName
	TestUserName2 = "test-user2"
	// TestPassword TestPassword
	TestPassword = "123456"
	// TestWrongPassword TestWrongPassword
	TestWrongPassword = "wrongpassword"
	// TestHashedPassword TestHashedPassword
	TestHashedPassword = "$2a$10$sSZRLWLaKu2JxPz9zpNjxek3N9UWMA82pyiEWoI1yXA.IE7KcMxTq"
	// TestToken TestToken
	TestToken = "test-token"
	// TestEmail TestEmail
	TestEmail = "test@example.com"
	// TestEmail2 TestEmail
	TestEmail2 = "test2@example.com"
	// TestFamilyID TestFamilyID
	TestFamilyID = uint64(1)
	// TestFamilyName TestFamilyName
	TestFamilyName = "test-family"
	// TestRoleHead TestRoleHead
	TestRoleHead = "head"
	// TestRoleMember TestRoleMember
	TestRoleMember = "member"
	// TestTaskID1 TestTaskID1
	TestTaskID1 = uint64(1)
	// TestTaskName1 TestTaskName1
	TestTaskName1 = "test-task1"
	// TestTaskStatusTodo TestTaskStatusTodo
	TestTaskStatusTodo = "todo"
	// TestTaskStatusDone TestTaskStatusDone
	TestTaskStatusDone = "done"
	// TestTaskDate TestTaskDate
	TestTaskDate = "2020-12-31"
	// TestSecretKey TestSecretKey
	TestSecretKey = "test_secret_key"
)

// RunTestMain runs tests with setups
func RunTestMain(m *testing.M) int {
	// テスト用の時間を設定する
	SetTestTime()
	defer ResetTime()
	SetGormTestTime()
	defer ResetGormTime()

	// テスト結果を返す
	return m.Run()
}

// SetTime sets global time
func SetTime(t time.Time) {
	NowFunc = func() time.Time { return t }
}

// SetTestTime sets global test time
func SetTestTime() {
	SetTime(GetTestTime())
}

// GetTestTime returns time of TestTime
func GetTestTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(2020, time.January, 1, 19, 0, 0, 0, loc)
}

// ResetTime resets global time
func ResetTime() {
	NowFunc = time.Now
}

// SetGormTestTime sets gorm test time
func SetGormTestTime() {
	gorm.NowFunc = func() time.Time { return GetGormTestTime() }
}

// GetGormTestTime returns time of TestTime
func GetGormTestTime() time.Time {
	return GetTestTime()
}

// ResetGormTime resets gorm time
func ResetGormTime() {
	gorm.NowFunc = time.Now
}

// AnyTime only asserts that argument is of time.Time type.
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// SetTestEnv sets temporarily environment variable for test
func SetTestEnv(key, val string) func() {
	preVal := os.Getenv(key)
	os.Setenv(key, val)
	return func() {
		os.Setenv(key, preVal)
	}
}
