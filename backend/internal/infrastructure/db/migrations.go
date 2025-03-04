package db

import (
	"app/internal/infrastructure/db/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (p *DBProvider) EnsureDatabase() error {
	// 1. Connect to MySQL without specifying DBName
	dsnNoDB := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		p.config.User,
		p.config.Password,
		p.config.Host,
		p.config.Port,
	)
	tmpDB, err := gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to MySQL (without DB): %v", err)
	}

	dbName := p.config.DBName

	// 2. Check if the database already exists
	var count int64
	checkSQL := "SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?"
	if err := tmpDB.Raw(checkSQL, dbName).Scan(&count).Error; err != nil {
		return fmt.Errorf("error checking database existence %s: %v", dbName, err)
	}

	if count > 0 {
		// 3a. If count>0, the database already exists
		log.Printf("⚠️ Database '%s' was already found, creation skipped.", dbName)
		return nil
	}

	// 3b. If the database does not exist, create it
	createSQL := fmt.Sprintf("CREATE DATABASE `%s`", dbName)
	if err := tmpDB.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("error creating database %s: %v", dbName, err)
	}

	log.Printf("✅ Database '%s' was created (it did not exist).", dbName)
	return nil
}

func Migrate(db *gorm.DB, creationDefaults bool) error {
	// 1. Execute AutoMigrate for needed entities
	if err := db.AutoMigrate(
		&models.ProviderModel{},
		&models.UserModel{},
		&models.ProfileModel{},
		&models.RoleModel{},
		&models.RefRoleUserModel{},
	); err != nil {
		return fmt.Errorf("autoMigrate error: %w", err)
	}

	// 2. (Optional) initialize default data, if you want
	if creationDefaults {
		if err := init_default_data(db); err != nil {
			return err
		}
	}

	log.Println("✅ Migration completed successfully.")
	return nil
}
