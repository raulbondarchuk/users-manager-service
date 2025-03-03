package models

import "app/internal/domain/user"

// RoleModel - model for storing roles
type RoleModel struct {
	ID   uint   `gorm:"primaryKey;column:id"`
	Role string `gorm:"column:role;size:255;not null;unique"`
	Desc string `gorm:"column:desc;size:255"`
}

func (RoleModel) TableName() string {
	return "roles"
}

// ToDomain - convert database model to domain structure
func (rm *RoleModel) ToDomain() *user.Role {
	return &user.Role{
		ID:   rm.ID,
		Role: rm.Role,
		Desc: rm.Desc,
	}
}
