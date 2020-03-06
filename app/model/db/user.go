package db

import "time"

// User represents a table
type User struct {
	ID        uint64 `gorm:"primary_key,column:id"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName is required by gorm
func (User) TableName() string {
	return "users"
}
