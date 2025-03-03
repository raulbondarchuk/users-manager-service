package db

import (
	"app/internal/infrastructure/db/models"
	"app/pkg/config"
	"log"

	"github.com/spf13/viper"
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

func init_Roles(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.RoleModel{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		if err := createRole(
			db,
			uint(viper.GetInt("database.migrations.defaults.roles.company.id")),
			viper.GetString("database.migrations.defaults.roles.company.name"),
			viper.GetString("database.migrations.defaults.roles.company.desc"),
		); err != nil {
			return err
		}

		if err := createRole(
			db,
			uint(viper.GetInt("database.migrations.defaults.roles.verificaciones.id")),
			viper.GetString("database.migrations.defaults.roles.verificaciones.name"),
			viper.GetString("database.migrations.defaults.roles.verificaciones.desc"),
		); err != nil {
			return err
		}

		if err := createRole(
			db,
			uint(viper.GetInt("database.migrations.defaults.roles.secondary.id")),
			viper.GetString("database.migrations.defaults.roles.secondary.name"),
			viper.GetString("database.migrations.defaults.roles.secondary.desc"),
		); err != nil {
			return err
		}

		if err := createRole(
			db,
			uint(viper.GetInt("database.migrations.defaults.roles.sat01.id")),
			viper.GetString("database.migrations.defaults.roles.sat01.name"),
			viper.GetString("database.migrations.defaults.roles.sat01.desc"),
		); err != nil {
			return err
		}
	}

	return nil
}

// Initialize default provider entity
func init_Provider(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.ProviderModel{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		if err := createProvider(
			db,
			uint(viper.GetInt("database.migrations.defaults.provider.liftel.id")),
			viper.GetString("database.migrations.defaults.provider.liftel.name"),
			viper.GetString("database.migrations.defaults.provider.liftel.desc"),
		); err != nil {
			return err
		}

		if err := createProvider(
			db,
			uint(viper.GetInt("database.migrations.defaults.provider.verificaciones.id")),
			viper.GetString("database.migrations.defaults.provider.verificaciones.name"),
			viper.GetString("database.migrations.defaults.provider.verificaciones.desc"),
		); err != nil {
			return err
		}

		if err := createProvider(
			db,
			uint(viper.GetInt("database.migrations.defaults.provider.secondary.id")),
			viper.GetString("database.migrations.defaults.provider.secondary.name"),
			viper.GetString("database.migrations.defaults.provider.secondary.desc"),
		); err != nil {
			return err
		}
	}

	return nil
}

// Initialize default user entity
func init_User(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.UserModel{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		// Create default user
		defaultUser := models.UserModel{
			ID:       uint(viper.GetInt("database.migrations.defaults.user.id")),
			Login:    config.ENV().DEFAULT_USER_LOGIN,
			Password: &config.ENV().DEFAULT_USER_PASSWORD,
		}
		if err := db.Create(&defaultUser).Error; err != nil {
			return err
		}
		log.Printf("Created default user: %s with ID: %d", defaultUser.Login, defaultUser.ID)

		// Assign company role to the default user
		if err := assignCompanyRoleToUser(db, defaultUser.ID); err != nil {
			return err
		}
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
