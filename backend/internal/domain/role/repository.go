package role

type RoleRepository interface {
	GetAllRoles() ([]Role, error)
	GetRoleByID(id uint) (*Role, error)
	AssignRoleToUser(userID, roleID uint) error
	GetUserRoles(userID uint) ([]Role, error)
}
