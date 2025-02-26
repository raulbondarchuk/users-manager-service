package env

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

// Config — struct, storing all needed environment variables
type Config struct {
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT" envDefault:"3307"`
	// More fields if needed
}

// variables for lazy initialization (singleton)
var (
	configInstance *Config
	once           sync.Once
)

// GetConfig — public access point to singleton
func Get() *Config {
	once.Do(func() {
		// 1. Try to load .env
		if err := godotenv.Load(); err != nil {
			log.Fatalf(" .env file not found, using system variables: %v", err)
		}

		// 2. Parse environment variables to Config struct
		cfg := Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("Error reading environment variables: %v", err)
		}
		configInstance = &cfg
	})
	return configInstance
}
