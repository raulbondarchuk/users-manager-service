package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// once_yaml — for lazy initialization (singleton)
var once_yaml sync.Once

// YAML — initialize Viper, reading YAML config.
// Can be called in main or inside composition/boot.
func YAML(configPath string) {
	once_yaml.Do(func() {
		// Set the path to the .yaml file
		viper.SetConfigFile(configPath)

		// Allow reading ENV variables (if needed)
		viper.AutomaticEnv()

		// Read the config
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading file %s: %v\n", configPath, err)
		} else {
			log.Printf("✅ Config loaded from %s\n", configPath)
		}
	})
}
