package db

import (
	"app/internal/infrastructure/db/models"
	"app/pkg/config"
	"log"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init_Roles(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.RoleModel{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		// Company role
		if err := createRole(
			db,
			uint(viper.GetInt("database.migrations.defaults.roles.company.id")),
			viper.GetString("database.migrations.defaults.roles.company.name"),
			viper.GetString("database.migrations.defaults.roles.company.desc"),
		); err != nil {
			return err
		}

		// Liftplay role
		if err := createRole(
			db,
			uint(viper.GetInt("database.migrations.defaults.roles.liftplay.id")),
			viper.GetString("database.migrations.defaults.roles.liftplay.name"),
			viper.GetString("database.migrations.defaults.roles.liftplay.desc"),
		); err != nil {
			return err
		}

		// Sat01 role
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
		// Liftel provider
		if err := createProvider(
			db,
			uint(viper.GetInt("database.migrations.defaults.provider.liftel.id")),
			viper.GetString("database.migrations.defaults.provider.liftel.name"),
			viper.GetString("database.migrations.defaults.provider.liftel.desc"),
		); err != nil {
			return err
		}

		// Verificaciones provider
		if err := createProvider(
			db,
			uint(viper.GetInt("database.migrations.defaults.provider.verificaciones.id")),
			viper.GetString("database.migrations.defaults.provider.verificaciones.name"),
			viper.GetString("database.migrations.defaults.provider.verificaciones.desc"),
		); err != nil {
			return err
		}

		// Secondary provider
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
		// Create default user ADMIN
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
