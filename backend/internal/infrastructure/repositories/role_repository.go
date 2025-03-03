package repositories

import (
	"app/internal/domain/user"
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/db/models"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() user.RoleRepository {
	return &roleRepository{db: db.GetProvider().GetDB()}
}

// GetAllRoles - получить все роли
func (r *roleRepository) GetAllRoles() ([]user.Role, error) {
	var roleModels []models.RoleModel
	if err := r.db.Find(&roleModels).Error; err != nil {
		return nil, err
	}

	// Используем ToDomain для конвертации
	roles := make([]user.Role, len(roleModels))
	for i, rm := range roleModels {
		roles[i] = *rm.ToDomain()
	}
	return roles, nil
}

// GetRoleByID - получить роль по ID
func (r *roleRepository) GetRoleByID(id uint) (*user.Role, error) {
	var roleModel models.RoleModel
	if err := r.db.First(&roleModel, id).Error; err != nil {
		return nil, err
	}
	return roleModel.ToDomain(), nil
}

// AssignRoleToUser - привязка роли к пользователю
func (r *roleRepository) AssignRoleToUser(userID, roleID uint) error {
	ref := models.RefRoleUserModel{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(&ref).Error
}

// GetUserRoles - получить все роли пользователя
func (r *roleRepository) GetUserRoles(userID uint) ([]user.Role, error) {
	var roleModels []models.RoleModel

	if err := r.db.Joins("JOIN ref_user_role ON ref_user_role.role_id = roles.id").
		Where("ref_user_role.user_id = ?", userID).
		Find(&roleModels).Error; err != nil {
		return nil, err
	}

	// Преобразуем в domain-модель
	roles := make([]user.Role, len(roleModels))
	for i, rm := range roleModels {
		roles[i] = *rm.ToDomain()
	}
	return roles, nil
}
