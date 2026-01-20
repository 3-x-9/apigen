
package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	BaseURL string
	Output string
	Timeout time.Duration

}

func Load(appName string, env string) *Config {
	_ = godotenv.Load()
	v := viper.New()

		serverMap := map[string]string{"http://localhost:8888": "http://localhost:8888"}

	v.SetDefault("base_url", "http://localhost:8888")
	v.SetDefault("output", "json")
	v.SetDefault("timeout", "30s")


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

	baseUrl := v.GetString("base_url")
	if env != "" {
		if url, ok := serverMap[strings.ToLower(env)]; ok {
			baseUrl = url
		}
	}

	timeout, _ := time.ParseDuration(v.GetString("timeout"))
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Config{
		BaseURL: baseUrl,
		Output: v.GetString("output"),
		Timeout: timeout,

	}
}
