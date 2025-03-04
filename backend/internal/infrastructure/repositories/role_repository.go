package repositories

import (
	"app/internal/domain/role"
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/db/models"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() role.RoleRepository {
	return &roleRepository{db: db.GetProvider().GetDB()}
}

// CreateRole - create new role
func (r *roleRepository) CreateRole(role *role.Role) (uint, error) {
	roleModel := models.RoleModel{
		Role: role.Role,
		Desc: role.Desc,
	}

	if err := r.db.Create(&roleModel).Error; err != nil {
		return 0, err
	}

	return roleModel.ID, nil
}

// GetAllRoles - get all roles
func (r *roleRepository) GetAllRoles() ([]role.Role, error) {
	var roleModels []models.RoleModel
	if err := r.db.Find(&roleModels).Error; err != nil {
		return nil, err
	}

	// Use ToDomain for conversion
	roles := make([]role.Role, len(roleModels))
	for i, rm := range roleModels {
		roles[i] = *rm.ToDomain()
	}
	return roles, nil
}

// GetRoleByID - get role by ID
func (r *roleRepository) GetRoleByID(id uint) (*role.Role, error) {
	var roleModel models.RoleModel
	if err := r.db.First(&roleModel, id).Error; err != nil {
		return nil, err
	}
	return roleModel.ToDomain(), nil
}

// AssignRoleToUser - assign role to user
func (r *roleRepository) AssignRoleToUser(userID, roleID uint) error {
	ref := models.RefRoleUserModel{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&ref).Error
}

// GetUserRoles - get user roles
func (r *roleRepository) GetUserRoles(userID uint) ([]role.Role, error) {
	var roleModels []models.RoleModel

	if err := r.db.Joins("JOIN ref_user_role ON ref_user_role.role_id = roles.id").
		Where("ref_user_role.user_id = ?", userID).
		Find(&roleModels).Error; err != nil {
		return nil, err
	}

	// Convert to domain model
	roles := make([]role.Role, len(roleModels))
	for i, rm := range roleModels {
		roles[i] = *rm.ToDomain()
	}
	return roles, nil
}

// RemoveRoleFromUser - remove role from user
func (r *roleRepository) RemoveRoleFromUser(userID, roleID uint) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.RefRoleUserModel{}).Error
}

// GetRoleByName - get role by name
func (r *roleRepository) GetRoleByName(name string) (*role.Role, error) {
	var roleModel models.RoleModel
	if err := r.db.Where("role = ?", name).First(&roleModel).Error; err != nil {
		return nil, err
	}
	return roleModel.ToDomain(), nil
}

// IsNotFoundError - check if error is a not found error
func (r *roleRepository) IsNotFoundError(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// GetRoleByNameWithTransaction gets a role by name within a transaction
func (r *roleRepository) GetRoleByNameWithTransaction(tx *gorm.DB, roleName string) (*role.Role, error) {
	var roleModel models.RoleModel
	if err := tx.Where("role = ?", roleName).First(&roleModel).Error; err != nil {
		return nil, err
	}
	return roleModel.ToDomain(), nil
}

// AssignRoleToUserWithTransaction assigns a role to a user within a transaction
func (r *roleRepository) AssignRoleToUserWithTransaction(tx *gorm.DB, userID, roleID uint) error {
	ref := models.RefRoleUserModel{
		UserID: userID,
		RoleID: roleID,
	}
	return tx.Create(&ref).Error
}
