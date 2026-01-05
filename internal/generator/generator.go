package generator

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (g *Generator) Generate(specPath, outputDir string, moduleName string) error {
	// 1. Load OpenAPI spec
	loader := openapi3.NewLoader()

	var doc *openapi3.T
	var err error

	if isURL(specPath) {
		parsedURL, err := url.Parse(specPath)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}
		doc, err = loader.LoadFromURI(parsedURL)
	} else {
		doc, err = loader.LoadFromFile(specPath)
	}

	if err != nil {
		return fmt.Errorf("failed to load spec: %w", err)
	}

	if err := createFolders(outputDir); err != nil {
		return err
	}

	// 2. Write go.mod
	if err := writeGoMod(outputDir, moduleName); err != nil {
		return err
	}

	// 3. Write config and root cmd
	if err := writeConfig(outputDir); err != nil {
		return err
	}
	if err := writeRootCmd(outputDir); err != nil {
		return err
	}

	// 4. Iterate paths and operations
	var cmdNames []string
	for path, pathItem := range doc.Paths.Map() {
		ops := map[string]*openapi3.Operation{
			"get":    pathItem.Get,
			"post":   pathItem.Post,
			"put":    pathItem.Put,
			"delete": pathItem.Delete,
		}

		for method, op := range ops {
			if op == nil {
				continue
			}
			cmdName := sanitizeCommandName(path, method)
			if err := writeEndpointCmd(outputDir, moduleName, cmdName, op, path, method); err != nil {
				return err
			}
			cmdNames = append(cmdNames, cmdName)
		}
	}

	if err := writeMain(outputDir, moduleName, cmdNames); err != nil {
		return err
	}

	fmt.Printf("Generated CLI in %s\n", outputDir)
	return nil
}

func createFolders(outputDir string) error {
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "cmd"),
		filepath.Join(outputDir, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

	}
	return nil
}

func writeGoMod(outputDir, moduleName string) error {
	goModContent := fmt.Sprintf(`module %s

go 1.22

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.19.0
)
`, moduleName)

	path := filepath.Join(outputDir, "go.mod")
	return os.WriteFile(path, []byte(goModContent), 0644)
}

func writeConfig(outputDir string) error {
	configCode := `
package config

import (
	"fmt"
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

func writeRootCmd(outputDir string) error {
	rootCode := `
	package cmd

	import (
		"github.com/spf13/cobra"
	)

	func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI is a command-line tool to interact with the API",
		}
		return cmd
		}
`
	path := filepath.Join(outputDir, "cmd", "root.go")
	return os.WriteFile(path, []byte(rootCode), 0644)
}

func sanitizeCommandName(path, method string) string {
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	return strings.Title(method + "_" + path)
}

func writeEndpointCmd(outputDir string, moduleName string, cmdName string, op *openapi3.Operation, path, method string) error {
	cmdCode := fmt.Sprintf(`
	package cmd

	import (
		"fmt"
		"github.com/spf13/cobra"
		"%s/config"
		"net/http"
		"io/ioutil"
		"encoding/json"
	)
	func New%sCmd() *cobra.Command {
		var limit int
		
		cmd := &cobra.Command{
			Use:   "%s",
			Short: "%s",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg := config.Load("%s")
				url := fmt.Sprintf("%%s%s?limit=%%d", cfg.BaseURL, limit)
				resp, err := http.Get(url)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				var pretty map[string]interface{}
				if err := json.Unmarshal(body, &pretty); err != nil {
					return err
				}
				prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(prettyJSON))
				return nil
			},	
		}
		cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of items")
		return cmd
	
		}	
`, moduleName, cmdName, cmdName, op.Summary, moduleName, path)

	pathFile := filepath.Join(outputDir, "cmd", strings.ToLower(cmdName)+".go")
	return os.WriteFile(pathFile, []byte(cmdCode), 0644)
}

func writeMain(outputDir string, moduleName string, cmds []string) error {
	cmdsInit := ""
	for _, c := range cmds {
		cmdsInit += fmt.Sprintf("\trootCmd.AddCommand(cmd.New%sCmd())\n", c)
	}

	mainCode := fmt.Sprintf(`
package main

import (
	"%s/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
%s
	cobra.CheckErr(rootCmd.Execute())
}
`, moduleName, cmdsInit)

	path := filepath.Join(outputDir, "main.go")
	return os.WriteFile(path, []byte(mainCode), 0644)
}
