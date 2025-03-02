package composition

import (
	"app/internal/infrastructure/db"
	http "app/internal/infrastructure/transport/http/server"
	"app/pkg/config"

	"github.com/spf13/viper"
)

const (
	YAML_PATH = "./configs/config.yaml"
)

// Helper functions for starting HTTP and gRPC servers, configs and database.
// Вспомогательные функции для запуска HTTP и gRPC серверов, конфигураций и базы данных.

func config_init() {
	config.INIT(YAML_PATH)
}

func db_init() {
	cfg := db.Config{}
	cfg.Set(
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		config.ENV().DBUser,
		config.ENV().DBPassword,
	)
	cfg.SetDBName(viper.GetString("database.schema"))
	cfg.SetEnsureDB(viper.GetBool("database.ensure"))
	cfg.SetAutoMigrate(viper.GetBool("database.auto_migrate"))
	cfg.SetCreationDefaults(viper.GetBool("database.migrations.creation_defaults"))
	db.NewDBProvider(cfg).Load()
}

func http_init() {
	http.MustLoad()
}
