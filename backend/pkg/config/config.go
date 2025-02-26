package config

func INIT(yamlPath string) {
	ENV()
	YAML(yamlPath)
}
