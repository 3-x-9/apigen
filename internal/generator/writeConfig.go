package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func writeConfig(outputDir string, schemes map[string]AuthScheme, servers openapi3.Servers) error {
	var structFields strings.Builder
	var viperDefaults strings.Builder
	var structAssigns strings.Builder

	for _, s := range schemes {
		// Clean name for struct field (simple sanitization)
		fieldName := s.Name
		if len(fieldName) > 0 {
			fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]
		}
		// e.g. "PetstoreAuth"
		fieldName += "Auth"

		// e.g. "PETSTORE_AUTH"
		envName := strings.ReplaceAll(strings.ToLower(s.Name), "-", "_") + "_auth"

		structFields.WriteString(fmt.Sprintf("\t%s string\n", fieldName))
		viperDefaults.WriteString(fmt.Sprintf("\tv.SetDefault(\"%s\", \"\")\n", envName))
		structAssigns.WriteString(fmt.Sprintf("\t\t%s: v.GetString(\"%s\"),\n", fieldName, envName))
	}

	serverMapCode := ""
	defaultBaseURL := "http://localhost:8080"
	if len(servers) > 0 {
		defaultBaseURL = servers[0].URL
		serverMapEntries := []string{}
		for _, s := range servers {
			name := strings.ToLower(s.Description)
			if name == "" {
				name = s.URL
			}
			// replace spaces with underscores/hyphens for easier CLI use
			name = strings.ReplaceAll(name, " ", "-")
			serverMapEntries = append(serverMapEntries, fmt.Sprintf("\"%s\": \"%s\"", name, s.URL))
		}
		serverMapCode = fmt.Sprintf("\tserverMap := map[string]string{%s}\n", strings.Join(serverMapEntries, ", "))
	}

	configCode := fmt.Sprintf(`
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
%s
}

func Load(appName string, env string) *Config {
	_ = godotenv.Load()
	v := viper.New()

	%s
	v.SetDefault("base_url", "%s")
	v.SetDefault("output", "json")
	v.SetDefault("timeout", "30s")
%s

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
%s
	}
}
`, structFields.String(), serverMapCode, defaultBaseURL, viperDefaults.String(), structAssigns.String())

	path := filepath.Join(outputDir, "config", "config.go")
	return os.WriteFile(path, []byte(configCode), 0644)
}
