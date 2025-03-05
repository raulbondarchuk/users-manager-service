package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

// once_yaml — for lazy initialization (singleton)
var once_yaml sync.Once

// YAML — initialize Viper, reading YAML config.
// Can be called in main or inside composition/boot.
func YAML(configPath string) {
	once_yaml.Do(func() {
		// Check if the specified config file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Printf("Config file %s not found, searching for any .yaml file in the root directory...\n", configPath)

			// Search for any .yaml file in the root directory
			rootFiles, err := filepath.Glob("*.yaml")
			if err != nil || len(rootFiles) == 0 {
				log.Fatalf("No .yaml config file found in the root directory: %v\n", err)
			}

			// Use the first .yaml file found
			configPath = rootFiles[0]
			log.Printf("Using config file %s found in the root directory\n", configPath)
		}

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
