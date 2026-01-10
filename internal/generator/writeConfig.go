package generator

import (
	"os"
	"path/filepath"
)

func writeConfig(outputDir string) error {
	configCode := `
package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string
	ApiKey string
	Output string
}

func Load(appName string) *Config {
	v := viper.New()

	v.SetDefault("base_url", "http://localhost:8080")
	v.SetDefault("output", "json")

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		v.AddConfigPath(filepath.Join(homeDir, "." + appName))
	}
	v.AddConfigPath(".")

	v.SetEnvPrefix(strings.ToUpper(appName))
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	return &Config{
		BaseURL: v.GetString("base_url"),
		ApiKey: v.GetString("api_key"),
		Output: v.GetString("output"),
	}
}
`

	path := filepath.Join(outputDir, "config", "config.go")
	return os.WriteFile(path, []byte(configCode), 0644)
}
