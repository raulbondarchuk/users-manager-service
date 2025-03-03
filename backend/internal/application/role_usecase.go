package application

import (
	"app/internal/domain/role"
	"app/internal/domain/user"
	"fmt"
)

// RoleUseCase - structure for processing business logic of roles
type RoleUseCase struct {
	roleRepo role.RoleRepository
	userRepo user.Repository
}

// NewRoleUseCase - constructor
func NewRoleUseCase(roleRepo role.RoleRepository, userRepo user.Repository) *RoleUseCase {
	return &RoleUseCase{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// GetAllRoles - get all roles
func (uc *RoleUseCase) GetAllRoles() ([]role.Role, error) {
	return uc.roleRepo.GetAllRoles()
}

// GetRoleByID - get role by ID
func (uc *RoleUseCase) GetRoleByID(id uint) (*role.Role, error) {
	return uc.roleRepo.GetRoleByID(id)
}

// GetRolesByUsername - get roles by username
func (uc *RoleUseCase) GetRolesByUsername(username string) ([]role.Role, error) {
	usr, err := uc.userRepo.GetByLogin(username)
	if err != nil {
		if uc.userRepo.IsNotFoundError(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return uc.roleRepo.GetUserRoles(usr.ID)
}
