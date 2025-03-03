package models

// RefRoleUserModel
type RefRoleUserModel struct {
	UserID uint `gorm:"primaryKey;column:user_id"`
	RoleID uint `gorm:"primaryKey;column:role_id"`

	User UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Role RoleModel `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
}

func (RefRoleUserModel) TableName() string {
	return "ref_user_role"
}
