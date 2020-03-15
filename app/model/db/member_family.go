package db

import "time"

const (
	// FamilyRoleHead the role in family
	FamilyRoleHead = "head"
	// FamilyRoleMember the role in family
	FamilyRoleMember = "member"
)

// MemberFamily represents a table
type MemberFamily struct {
	MemberID  uint64 `gorm:"primary_key,column:member_id"`
	FamilyID  uint64 `gorm:"primary_key,column:family_id"`
	Role      string `gorm:"column:role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName is required by gorm
func (MemberFamily) TableName() string {
	return "members_families"
}
