package models

import "app/internal/domain/role"

type RefRoleUserModel struct {
	UserID uint `gorm:"column:userId;primaryKey"`
	RoleID uint `gorm:"column:roleId;primaryKey"`
}

func (RefRoleUserModel) TableName() string {
	return "ref_user_role"
}

// ToDomain - convert to domain structure
func (rru *RefRoleUserModel) ToDomain() *role.RefRoleUser {
	return &role.RefRoleUser{
		UserID: rru.UserID,
		RoleID: rru.RoleID,
	}
}
