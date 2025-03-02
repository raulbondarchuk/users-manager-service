package db

import (
	logger_gorm "app/pkg/logger/gorm"
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	CustomLogger     bool // if true => use custom logger
}

func (c *Config) Set(host string, port int, user, password string) {
	c.Host = host
	c.Port = port
	c.User = user
	c.Password = password
}

// Method for setting the DBName field. Is the name of the schema to connect to.
func (c *Config) SetDBName(dbname string) { c.DBName = dbname }

// Method for setting the EnsureDB flag.
// If the flag is true, the database will be created if it does not exist.
func (c *Config) SetEnsureDB(ensure bool) { c.EnsureDB = ensure }

// Method for setting the AutoMigrate flag. If the flag is true, the migrations will be executed.
func (c *Config) SetAutoMigrate(autoMigrate bool) { c.AutoMigrate = autoMigrate }

// Method for setting the CreationDefaults flag. If the flag is true, the creation of default entities will be executed.
func (c *Config) SetCreationDefaults(creationDefaults bool) { c.CreationDefaults = creationDefaults }

// Method for setting the CustomLogger flag. If the flag is true, the custom logger will be used.
func (c *Config) SetCustomLogger(customLogger bool) { c.CustomLogger = customLogger }

type DBProvider struct {
	config Config
	db     *gorm.DB
	mu     sync.Mutex // For thread safety during initialization
}

var (
	globalProvider *DBProvider
	once           sync.Once
)

// Initialize creates and saves the provider in the global variable
func Initialize(cfg Config) {
	var initErr error
	once.Do(func() {
		globalProvider = NewDBProvider(cfg)
		initErr = globalProvider.Load()
		if initErr != nil {
			log.Fatal(initErr)
		}
	})
}

// GetProvider returns the global provider
func GetProvider() *DBProvider {
	if globalProvider == nil {
		log.Fatal("Database provider not initialized. Call Initialize() first")
	}
	return globalProvider
}

// getLogger initializes and returns the appropriate logger based on the configuration.
func getLogger(config *Config) logger.Interface {
	if config.CustomLogger {
		return logger_gorm.NewCustomLogger()
	}
	return logger.Default
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: getLogger(&p.config), // Use the logger from the new method
	})
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
func (p *DBProvider) GetDB() *gorm.DB {
	if p.db == nil {
		p.Load()
	}
	return p.db
}
