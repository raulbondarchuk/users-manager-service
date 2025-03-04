package role

import "gorm.io/gorm"

type RoleRepository interface {
	GetAllRoles() ([]Role, error)
	GetRoleByID(id uint) (*Role, error)
	AssignRoleToUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]Role, error)
	CreateRole(role *Role) (uint, error)
	RemoveRoleFromUser(userID, roleID uint) error
	GetRoleByName(name string) (*Role, error)
	IsNotFoundError(err error) bool

	GetRoleByNameWithTransaction(tx *gorm.DB, roleName string) (*Role, error)
	AssignRoleToUserWithTransaction(tx *gorm.DB, userID, roleID uint) error
}
