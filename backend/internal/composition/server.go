package composition

import (
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/transport/http"
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
		config.ENV().DBHost,
		config.ENV().DBPort,
		config.ENV().DBUser,
		config.ENV().DBPassword,
	)
	cfg.SetDBName(viper.GetString("database.schema"))
	db.NewDBProvider(cfg).Load()
}

func http_init() {
	http.MustLoad()
}
