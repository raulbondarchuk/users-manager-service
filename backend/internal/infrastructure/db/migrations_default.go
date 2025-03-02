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

	if err := init_User(db); err != nil {
		return err
	}

	return nil
}

// Initialize default provider entity
func init_Provider(db *gorm.DB) error {
	// Example: Check if there is at least one
	var count int64
	if err := db.Model(&models.ProviderModel{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		// Create default provider
		defaultProvider := models.ProviderModel{
			ID:   uint(viper.GetInt("database.migrations.defaults.provider.id")),
			Name: viper.GetString("database.migrations.defaults.provider.name"),
			Desc: viper.GetString("database.migrations.defaults.provider.desc"),
		}
		if err := db.Create(&defaultProvider).Error; err != nil {
			return err
		}
		log.Printf("Created default provider: %s with ID: %d", defaultProvider.Name, defaultProvider.ID)
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
			Password: config.ENV().DEFAULT_USER_PASSWORD,
		}
		if err := db.Create(&defaultUser).Error; err != nil {
			return err
		}
		log.Printf("Created default user: %s with ID: %d", defaultUser.Login, defaultUser.ID)
	}
	return nil
}
