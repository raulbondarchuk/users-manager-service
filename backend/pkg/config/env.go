package config

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

// Config — struct, storing all needed environment variables
type Config struct {

	// DATABASE
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`

	// VERIFICACIONES
	VERIFICACIONES_USERNAME string `env:"VERIFICACIONES_USERNAME,required"`
	VERIFICACIONES_PASSWORD string `env:"VERIFICACIONES_PASSWORD,required"`

	// DEFAULT USER
	DEFAULT_USER_LOGIN    string `env:"DEFAULT_USER_LOGIN,required"`
	DEFAULT_USER_PASSWORD string `env:"DEFAULT_USER_PASSWORD,required"`

	// PASETO + REFRESH
	PASETO_SK               string `env:"PASETO_SK,required"`
	PASETO_EXPIRATION_TIME  string `env:"PASETO_EXPIRATION_TIME,required"`
	REFRESH_EXPIRATION_TIME string `env:"REFRESH_EXPIRATION_TIME,required"`

	// SMTP
	MAIL_SMTP_HOST     string `env:"MAIL_SMTP_HOST,required"`
	MAIL_SMTP_PORT     string `env:"MAIL_SMTP_PORT,required"`
	MAIL_SMTP_USERNAME string `env:"MAIL_SMTP_USERNAME,required"`
	MAIL_SMTP_PASSWORD string `env:"MAIL_SMTP_PASSWORD,required"`
	MAIL_SMTP_TLS      bool   `env:"MAIL_SMTP_TLS,required"`

	// MIDDLEWARE PARA CONTRASEÑAS
	MIDDLEWARE_PASSWORD string `env:"MIDDLEWARE_PASSWORD,required"`
	// More fields if needed ...
}

// variables for lazy initialization (singleton)
var (
	configInstance *Config
	once_env       sync.Once
)

// GetConfig — public access point to singleton
func ENV() *Config {
	once_env.Do(func() {
		// 1. Try to load .env
		if err := godotenv.Load(); err != nil {
			log.Fatalf(" .env file not found, using system variables: %v", err)
		}
		log.Printf("✅ .env file loaded")
		// 2. Parse environment variables to Config struct
		cfg := Config{}
		if err := env.Parse(&cfg); err != nil {
			log.Fatalf("Error reading environment variables: %v", err)
		}
		configInstance = &cfg
	})
	return configInstance
}
