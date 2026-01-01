package api

import (
	"github.com/gren-95/subCli/internal/config"
)

// Exported for backwards compatibility
var AppConfig = &config.AppConfig

func LoadConfig() error {
	return config.LoadConfig()
}

func SaveConfig() error {
	return config.SaveConfig()
}
