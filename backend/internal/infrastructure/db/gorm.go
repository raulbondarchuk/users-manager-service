package db

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func (c *Config) Set(host string, port int, user, password string) {
	c.Host = host
	c.Port = port
	c.User = user
	c.Password = password
}

func (c *Config) SetDBName(dbname string) {
	c.DBName = dbname
}

type DBProvider struct {
	config Config
	db     *gorm.DB
	mu     sync.Mutex // For thread safety during initialization
}

// NewDBProvider — constructor, which remembers the config (but does not open the connection immediately).
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

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.config.User, p.config.Password, p.config.Host, p.config.Port, p.config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	log.Printf("✅ Successfully connected to database %s at %s:%d", p.config.DBName, p.config.Host, p.config.Port)
	p.db = db
	return nil
}

// GetDB — method for returning the already initialized *gorm.DB
func (p *DBProvider) Get() *gorm.DB {
	if p.db == nil {
		p.Load()
	}
	return p.db

}

func Migrate(db *gorm.DB) error {
	// return db.AutoMigrate(
	// 	&user.User{},
	// 	&user.Profile{},
	// 	&subuser.SubUser{},
	// )
	return nil
}
