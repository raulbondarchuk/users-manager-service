package composition

import "app/pkg/config"

const (
	YAML_PATH = "./configs/config.yaml"
)

// Helper functions for starting HTTP and gRPC servers.
// Вспомогательные функции для запуска HTTP и gRPC серверов.

func config_init() {
	config.INIT(YAML_PATH)
}
