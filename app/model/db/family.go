package db

import "time"

// Family represents a table
type Family struct {
	ID        uint64 `gorm:"primary_key,column:id"`
	Name      string `gorm:"column:name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName is required by gorm
func (Family) TableName() string {
	return "families"
}
