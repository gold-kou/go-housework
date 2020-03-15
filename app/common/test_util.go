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
	TestUserID = 1
	// TestUserID2 TestUserID2
	TestUserID2 = 2
	// TestUserName TestUserName
	TestUserName = "test-user"
	// TestPassword TestPassword
	TestPassword = "123456"
	// TestWrongPassword TestWrongPassword
	TestWrongPassword = "wrongpassword"
	// TestHashedPassword TestHashedPassword
	TestHashedPassword = "$2a$10$sSZRLWLaKu2JxPz9zpNjxek3N9UWMA82pyiEWoI1yXA.IE7KcMxTq"
	// TestEmail TestEmail
	TestEmail = "test@example.com"
	// TestFamilyID TestFamilyID
	TestFamilyID = 1
	// TestFamilyName TestFamilyName
	TestFamilyName = "test-family"
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
