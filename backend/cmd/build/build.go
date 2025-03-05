package build

import "github.com/raulbondarchuk/fast-go/builder"

func Build() {
	builderConfig := builder.BuildConfig{
		DefaultMode:      "dev", // Por ejemplo: local, dev, prod
		OutputFilename:   "um_v3",
		OutputDir:        "./",
		SourceFile:       "./cmd/main.go",
		BuildLinux:       true,
		BuildWindows:     false,
		PossibleDirs:     []string{"", "configs", "cfg", "config", "internal/config"},
		ConfigExtensions: []string{"toml", "yaml"},
		AddAppOnConfig:   false,
	}
	builderConfig.Run()
}
