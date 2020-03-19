package db

import "time"

// Task represents a table
type Task struct {
	ID        uint64 `gorm:"primary_key,column:id"`
	Name      string `gorm:"column:name"`
	MemberID  uint64 `gorm:"column:member_id"`
	FamilyID  uint64 `gorm:"column:family_id"`
	Status    string `gorm:"column:status"`
	Date      string `gorm:"column:date"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName is required by gorm
func (Task) TableName() string {
	return "tasks"
}
