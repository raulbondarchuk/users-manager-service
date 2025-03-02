package db

import (
	"app/internal/domain/provider"
	"log"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func init_default_data(db *gorm.DB) error {
	if err := init_Provider(db); err != nil {
		return err
	}
	return nil
}

// Initialize default provider entity
func init_Provider(db *gorm.DB) error {
	// Example: Check if there is at least one
	var count int64
	if err := db.Model(&provider.Provider{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		// Create default provider
		defaultProvider := provider.Provider{
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
