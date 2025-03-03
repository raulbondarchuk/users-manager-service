package user

import (
	"fmt"
	"time"

	"app/internal/domain/role"
)

type UserService struct {
	userRepo Repository
	roleRepo role.RoleRepository
}

func NewUserService(userRepo Repository, roleRepo role.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// EnsureUserRoles checks if user has required roles and assigns them if not
func (s *UserService) EnsureUserRoles(usr *User) error {
	// Get user roles
	userRoles, err := s.roleRepo.GetUserRoles(usr.ID)
	if err != nil {
		return fmt.Errorf("get user roles error: %w", err)
	}

	// Check if user has required roles
	hasCompanyRole := false
	hasCompanyIDRole := false
	companyIDRoleName := fmt.Sprintf("company_%d", usr.CompanyID)

	for _, r := range userRoles {
		if r.Role == "company" {
			hasCompanyRole = true
		}
		if r.Role == companyIDRoleName {
			hasCompanyIDRole = true
		}
	}

	// Get all roles to check if required roles exist
	allRoles, err := s.roleRepo.GetAllRoles()
	if err != nil {
		return fmt.Errorf("get all roles error: %w", err)
	}

	// Check if roles exist in DB
	companyRoleID := uint(0)
	companyIDRoleID := uint(0)

	for _, r := range allRoles {
		if r.Role == "company" {
			companyRoleID = r.ID
		}
		if r.Role == companyIDRoleName {
			companyIDRoleID = r.ID
		}
	}

	// Create "company" role if it doesn't exist
	if companyRoleID == 0 {
		newRole := &role.Role{
			Role: "company",
			Desc: "General company role",
		}
		companyRoleID, err = s.roleRepo.CreateRole(newRole)
		if err != nil {
			return fmt.Errorf("create company role error: %w", err)
		}
	}

	// Create company_ID role if it doesn't exist
	if companyIDRoleID == 0 {
		newRole := &role.Role{
			Role: companyIDRoleName,
			Desc: fmt.Sprintf("Role for company ID %d", usr.CompanyID),
		}
		companyIDRoleID, err = s.roleRepo.CreateRole(newRole)
		if err != nil {
			return fmt.Errorf("create company ID role error: %w", err)
		}
	}

	// Determine if user should have company role
	// If UserOwner is set (not nil and not 0), don't assign company role
	shouldHaveCompanyRole := usr.OwnerID == nil || *usr.OwnerID == 0

	// Assign "company" role to user if needed and if user doesn't have UserOwner
	if !hasCompanyRole && shouldHaveCompanyRole {
		if err := s.roleRepo.AssignRoleToUser(usr.ID, companyRoleID); err != nil {
			return fmt.Errorf("assign company role error: %w", err)
		}
	}

	// Assign company_ID role to user if needed
	if !hasCompanyIDRole {
		if err := s.roleRepo.AssignRoleToUser(usr.ID, companyIDRoleID); err != nil {
			return fmt.Errorf("assign company ID role error: %w", err)
		}
	}

	return nil
}

// GetUserRoles gets user roles from database
func (s *UserService) GetUserRoles(userID uint) ([]role.Role, error) {
	userRoles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles error: %w", err)
	}

	// If roles are empty, try again after a short delay
	if len(userRoles) == 0 {
		time.Sleep(100 * time.Millisecond)
		userRoles, err = s.roleRepo.GetUserRoles(userID)
		if err != nil {
			return nil, fmt.Errorf("get user roles error after retry: %w", err)
		}
	}

	return userRoles, nil
}

// GetRoleNamesString converts roles to comma-separated string
func (s *UserService) GetRoleNamesString(roles []role.Role) string {
	roleNames := ""
	for i, r := range roles {
		if i > 0 {
			roleNames += ","
		}
		roleNames += r.Role
	}
	return roleNames
}
