package models

import "app/internal/domain/role"

// RoleModel - model for storing roles
type RoleModel struct {
	ID   uint   `gorm:"primaryKey;column:id"`
	Role string `gorm:"column:role;size:255;not null;unique"`
	Desc string `gorm:"column:desc;size:255"`

	// Many-to-many relationship through ref_user_role
	Roles []RoleModel `gorm:"many2many:ref_user_role;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:role_id"`
}

func (RoleModel) TableName() string {
	return "roles"
}

// ToDomain - convert database model to domain structure
func (rm *RoleModel) ToDomain() *role.Role {
	return &role.Role{
		ID:   rm.ID,
		Role: rm.Role,
		Desc: rm.Desc,
	}
}
