package common

import (
	"database/sql/driver"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

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
	return time.Date(2018, time.May, 21, 23, 0, 0, 0, loc)
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

// IntToPtr returns pointer of int
func IntToPtr(i int) *int {
	return &i
}

// SetTestEnv sets temporarily environment variable for test
func SetTestEnv(key, val string) func() {
	preVal := os.Getenv(key)
	os.Setenv(key, val)
	return func() {
		os.Setenv(key, preVal)
	}
}

// AnyTime only asserts that argument is of time.Time type.
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
