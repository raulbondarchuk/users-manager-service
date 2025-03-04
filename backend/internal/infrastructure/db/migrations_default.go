package db

import (
	"app/internal/infrastructure/db/models"
	"log"

	"gorm.io/gorm"
)

func init_default_data(db *gorm.DB) error {

	if err := init_Provider(db); err != nil {
		return err
	}

	if err := init_Roles(db); err != nil {
		return err
	}

	if err := init_User(db); err != nil {
		return err
	}

	return nil
}

func assignCompanyRoleToUser(db *gorm.DB, userID uint) error {
	var companyRole models.RoleModel
	if err := db.Where("role = ?", "company").First(&companyRole).Error; err != nil {
		return err
	}

	userRole := models.RefRoleUserModel{
		UserID: userID,
		RoleID: companyRole.ID,
	}
	if err := db.Create(&userRole).Error; err != nil {
		return err
	}
	log.Printf("Assigned company role to user ID: %d", userID)
	return nil
}

func createProvider(db *gorm.DB, id uint, name, desc string) error {

	provider := models.ProviderModel{
		ID:   id,
		Name: name,
		Desc: desc,
	}
	if err := db.Create(&provider).Error; err != nil {
		return err
	}
	log.Printf("Created default provider: %s with ID: %d", provider.Name, provider.ID)
	return nil
}

func createRole(db *gorm.DB, id uint, name, desc string) error {
	role := models.RoleModel{
		ID:   id,
		Role: name,
		Desc: desc,
	}
	if err := db.Create(&role).Error; err != nil {
		return err
	}
	log.Printf("Created default role: %s with ID: %d", role.Role, role.ID)
	return nil
}
