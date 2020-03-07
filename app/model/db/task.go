package db

import "time"

// Task represents a table
type Task struct {
	ID        uint64 `gorm:"primary_key,column:id"`
	Name      string `gorm:"column:name"`
	MemberID  uint64 `gorm:"column:member_id"`
	Status    string `gorm:"column:status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName is required by gorm
func (Task) TableName() string {
	return "tasks"
}
