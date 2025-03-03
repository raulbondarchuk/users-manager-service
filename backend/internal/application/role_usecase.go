package application

import (
	"app/internal/domain/role"
	"app/internal/domain/user"
	"fmt"
	"strings"
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

// AssignRolesToUser - assign roles to user by username
func (uc *RoleUseCase) AssignRolesToUser(username string, roleNames string) error {
	// 1. Get user by username
	usr, err := uc.userRepo.GetByLogin(username)
	if err != nil {
		if uc.userRepo.IsNotFoundError(err) {
			return fmt.Errorf("user not found: %s", username)
		}
		return fmt.Errorf("error retrieving user: %w", err)
	}

	// 2. Get existing user roles
	existingRoles, err := uc.roleRepo.GetUserRoles(usr.ID)
	if err != nil {
		return fmt.Errorf("error getting user roles: %w", err)
	}

	// Create a map for quick lookup of existing roles
	existingRoleMap := make(map[string]bool)
	for _, r := range existingRoles {
		existingRoleMap[r.Role] = true
	}

	// 3. Get all roles from the system
	allRoles, err := uc.roleRepo.GetAllRoles()
	if err != nil {
		return fmt.Errorf("error getting all roles: %w", err)
	}

	// Create a map for quick lookup of existing system roles
	systemRoleMap := make(map[string]uint)
	for _, r := range allRoles {
		systemRoleMap[r.Role] = r.ID
	}

	// 4. Process each role name from the input
	roleNamesSlice := strings.Split(roleNames, ",")
	for _, roleName := range roleNamesSlice {
		roleName = strings.TrimSpace(roleName)
		if roleName == "" {
			continue
		}

		// Check if user already has this role
		if existingRoleMap[roleName] {
			// User already has this role, skip
			continue
		}

		// Check if role exists in the system
		roleID, exists := systemRoleMap[roleName]
		if !exists {
			// Role doesn't exist, create it
			newRole := &role.Role{
				Role: roleName,
				Desc: fmt.Sprintf("Role %s", roleName),
			}
			var err error
			roleID, err = uc.roleRepo.CreateRole(newRole)
			if err != nil {
				return fmt.Errorf("error creating role %s: %w", roleName, err)
			}
		}

		// Assign role to user
		if err := uc.roleRepo.AssignRoleToUser(usr.ID, roleID); err != nil {
			return fmt.Errorf("error assigning role %s to user: %w", roleName, err)
		}
	}

	return nil
}

// EliminateRolesOfUser - remove roles from user by username
func (uc *RoleUseCase) EliminateRolesOfUser(username string, roleNames string) error {
	// 1. Get user by username
	usr, err := uc.userRepo.GetByLogin(username)
	if err != nil {
		if uc.userRepo.IsNotFoundError(err) {
			return fmt.Errorf("user not found: %s", username)
		}
		return fmt.Errorf("error retrieving user: %w", err)
	}

	// 2. Get existing user roles
	existingRoles, err := uc.roleRepo.GetUserRoles(usr.ID)
	if err != nil {
		return fmt.Errorf("error getting user roles: %w", err)
	}

	// Create a map for quick lookup of existing roles
	existingRoleMap := make(map[string]uint)
	for _, r := range existingRoles {
		existingRoleMap[r.Role] = r.ID
	}

	// 3. Process each role name from the input
	roleNamesSlice := strings.Split(roleNames, ",")
	for _, roleName := range roleNamesSlice {
		roleName = strings.TrimSpace(roleName)
		if roleName == "" {
			continue
		}

		// Check if user has this role
		roleID, hasRole := existingRoleMap[roleName]
		if !hasRole {
			// User doesn't have this role, skip
			continue
		}

		// Remove role from user
		if err := uc.roleRepo.RemoveRoleFromUser(usr.ID, roleID); err != nil {
			return fmt.Errorf("error removing role %s from user: %w", roleName, err)
		}
	}

	return nil
}
