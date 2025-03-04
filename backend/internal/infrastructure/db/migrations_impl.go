package db

import (
	"app/internal/infrastructure/db/models"
	"app/pkg/config"
	"fmt"
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

func init_InternalCompany(db *gorm.DB) error {
	// 1. Check if the internal_companies table is empty
	var count int64
	if err := db.Model(&models.InternalCompanyModel{}).Count(&count).Error; err != nil {
		return fmt.Errorf("error counting internal_companies: %w", err)
	}

	// 2. Set initial auto-increment value if the table is empty
	if count == 0 {
		if err := db.Exec("ALTER TABLE internal_companies AUTO_INCREMENT = 20000").Error; err != nil {
			return fmt.Errorf("error setting initial auto-increment value: %w", err)
		}
	}

	// 3. Check if the trigger exists and create it if it doesn't
	triggerExistsQuery := `
		SELECT COUNT(*)
		FROM information_schema.triggers
		WHERE trigger_name = 'before_insert_internal_companies'
		AND trigger_schema = DATABASE();
	`

	var triggerCount int64
	if err := db.Raw(triggerExistsQuery).Scan(&triggerCount).Error; err != nil {
		return fmt.Errorf("error checking trigger existence: %w", err)
	}

	if triggerCount == 0 {
		createTriggerQuery := `
			CREATE TRIGGER before_insert_internal_companies
			BEFORE INSERT ON internal_companies
			FOR EACH ROW
			BEGIN
				DECLARE max_id INT;

				-- Get the maximum ID from the table
				SELECT IFNULL(MAX(ID), 19999) INTO max_id FROM internal_companies;

				-- Set the new ID as max_id + 1
				SET NEW.ID = max_id + 1;
			END;
		`

		if err := db.Exec(createTriggerQuery).Error; err != nil {
			return fmt.Errorf("error creating trigger: %w", err)
		}
		log.Println("Trigger 'before_insert_internal_companies' created successfully.")
	}

	return nil
}
