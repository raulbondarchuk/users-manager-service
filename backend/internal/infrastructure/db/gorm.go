package db

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config — struct for storing the database configuration.
type Config struct {
	Host             string
	Port             int
	User             string
	Password         string
	DBName           string
	EnsureDB         bool // if true => execute CREATE DATABASE IF NOT EXISTS
	AutoMigrate      bool // if true => execute AutoMigrate
	CreationDefaults bool // if true => execute creation of default entities
}

func (c *Config) Set(host string, port int, user, password string) {
	c.Host = host
	c.Port = port
	c.User = user
	c.Password = password
}

// Method for setting the DBName field. Is the name of the schema to connect to.
func (c *Config) SetDBName(dbname string) {
	c.DBName = dbname
}

// Method for setting the EnsureDB flag.
// If the flag is true, the database will be created if it does not exist.
func (c *Config) SetEnsureDB(ensure bool) {
	c.EnsureDB = ensure
}

// Method for setting the AutoMigrate flag. If the flag is true, the migrations will be executed.
func (c *Config) SetAutoMigrate(autoMigrate bool) {
	c.AutoMigrate = autoMigrate
}

// Method for setting the CreationDefaults flag. If the flag is true, the creation of default entities will be executed.
func (c *Config) SetCreationDefaults(creationDefaults bool) {
	c.CreationDefaults = creationDefaults
}

type DBProvider struct {
	config Config
	db     *gorm.DB
	mu     sync.Mutex // For thread safety during initialization
}

// Constructor, which remembers the config (but does not open the connection immediately).
func NewDBProvider(cfg Config) *DBProvider {
	return &DBProvider{config: cfg}
}

// Load — opens the connection and (optionally) calls migrations.
// Call once when starting the application.
func (p *DBProvider) Load() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.db != nil {
		// already initialized
		return nil
	}

	// 1. If needed, create database (EnsureDB)
	if p.config.EnsureDB {
		if err := p.EnsureDatabase(); err != nil {
			return fmt.Errorf("error ensuring database: %w", err)
		}
	}

	// 2. Connect to the database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.config.User,
		p.config.Password,
		p.config.Host,
		p.config.Port,
		p.config.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Printf("✅ Successfully connected to database %s at %s:%d",
		p.config.DBName, p.config.Host, p.config.Port)
	p.db = db

	// 3. If needed, execute migration
	if p.config.AutoMigrate {
		if err := Migrate(p.db, p.config.CreationDefaults); err != nil {
			return fmt.Errorf("migration error: %v", err)
		}
	}

	return nil
}

// GetDB — method for returning the already initialized *gorm.DB
func (p *DBProvider) Get() *gorm.DB {
	if p.db == nil {
		p.Load()
	}
	return p.db

}
