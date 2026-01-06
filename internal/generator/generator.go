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

	authType, authHeader := detectAuth(doc)

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
			if err := writeEndpointCmd(outputDir, moduleName, cmdName, op, path, method, authType, authHeader); err != nil {
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

func detectAuth(doc *openapi3.T) (authType string, headerName string) {
	if doc == nil || doc.Components == nil || doc.Components.SecuritySchemes == nil {
		return "", ""
	}
	for name, schemeRef := range doc.Components.SecuritySchemes {
		if schemeRef == nil || schemeRef.Value == nil {
			continue
		}
		s := schemeRef.Value

		switch s.Type {
		case "http":
			if s.Scheme == "bearer" {
				return "bearer", "Authorization"
			}
		case "apiKey":
			if s.In == "header" {
				return "apiKey", s.Name
			}
		}
		_ = name
	}
	return "", ""
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

func writeEndpointCmd(outputDir string, moduleName string, cmdName string, op *openapi3.Operation, path, method string, authType string, authHeader string) error {
	// Build parameter-driven code pieces from the operation parameters
	var imports = map[string]bool{
		"fmt":                                true,
		"github.com/spf13/cobra":             true,
		fmt.Sprintf("%s/config", moduleName): true,
		"net/http":                           true,
		"io/ioutil":                          true,
		"encoding/json":                      true,
		"net/url":                            true,
		"strings":                            true,
	}

	var varDecls strings.Builder
	var flagsSetup strings.Builder
	var pathReplacements strings.Builder
	var queryBuild strings.Builder
	var needStrconv bool
	var authCode string

	if authType == "bearer" {
		authCode = `
		if cfg.ApiKey != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
		}
		`
	} else if authType == "apiKey" {
		authCode = fmt.Sprintf(`
		if cfg.ApiKey != "" {
			req.Header.Set("%s", cfg.ApiKey)}`, authHeader)
	}
	// default common param used in many generated commands
	varDecls.WriteString("\tvar limit int\n\n")
	flagsSetup.WriteString("\tcmd.Flags().IntVar(&limit, \"limit\", 10, \"Maximum number of items\")\n")

	// add body flag for POST and PUT methods
	var bodyHandling string
	var headerHandling string
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
		varDecls.WriteString("\tvar body string\n")
		flagsSetup.WriteString("\tcmd.Flags().StringVarP(&body, \"body\", \"b\", \"\", \"Request body (raw JSON, @filename, or '-' for stdin)\")\n")
		imports["io"] = true
		imports["os"] = true
		imports["bytes"] = true
		bodyHandling = `
			// prepare request body (allow raw JSON, @filename, or '-' for stdin)
			var bodyReader io.Reader
			if body != "" {
				if strings.HasPrefix(body, "@") {
					fname := strings.TrimPrefix(body, "@")
					var data []byte
					var err error
					if fname == "-" {
						data, err = ioutil.ReadAll(os.Stdin)
						if err != nil {
							return err
						}
					} else {
						data, err = ioutil.ReadFile(fname)
						if err != nil {
							return err
						}
					}
					bodyReader = bytes.NewReader(data)
				} else if body == "-" {
					data, err := ioutil.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					bodyReader = bytes.NewReader(data)
				} else {
					bodyReader = strings.NewReader(body)
				}
			}
`
		headerHandling = `
			if bodyReader != nil {
				req.Header.Set("Content-Type", "application/json")
			}
`
	} else if strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "DELETE" {
		// no body for GET/DELETE
		imports["io"] = true
		bodyHandling = `
			var bodyReader io.Reader = nil
`
		headerHandling = ""
	}

	for _, pRef := range op.Parameters {
		if pRef == nil || pRef.Value == nil {
			continue
		}
		p := pRef.Value
		name := p.Name
		in := p.In

		// determine type (default to string)
		goType := "string"
		flagFunc := "StringVar"
		// varName used in generated code
		varName := sanitizeVar(name)

		if p.Schema != nil && p.Schema.Value != nil {
			types := p.Schema.Value.Type
			if len(*types) > 0 {
				switch (*types)[0] {
				case "string":
					goType = "string"
					flagFunc = "StringVar"
				case "integer":
					goType = "int"
					flagFunc = "IntVar"
					imports["strconv"] = true
				case "boolean":
					goType = "bool"
					flagFunc = "BoolVar"
					imports["strconv"] = true
				default:
					// treat other types as string for MVP
					goType = "string"
					flagFunc = "StringVar"
				}
			}
		}

		varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))

		// flag registration
		defaultVal := "\"\""
		if goType == "int" {
			defaultVal = "0"
		} else if goType == "bool" {
			defaultVal = "false"
		}
		flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s %s parameter\")\n", flagFunc, varName, name, defaultVal, in, name))

		// code to handle where the param goes
		if in == "path" {
			// replace path placeholders with param values
			replacement := varName
			if goType == "int" {
				replacement = fmt.Sprintf("strconv.Itoa(%s)", varName)
				needStrconv = true
			} else if goType == "bool" {
				replacement = fmt.Sprintf("strconv.FormatBool(%s)", varName)
				needStrconv = true
			}
			pathReplacements.WriteString(fmt.Sprintf("\tpathWithParams = strings.ReplaceAll(pathWithParams, \"{%s}\", %s)\n", name, replacement))
		} else if in == "query" {
			if goType == "string" {
				queryBuild.WriteString(fmt.Sprintf("\tif %s != \"\" { q.Set(\"%s\", %s) }\n", varName, name, varName))
			} else if goType == "int" {
				replacement := fmt.Sprintf("strconv.Itoa(%s)", varName)
				queryBuild.WriteString(fmt.Sprintf("\tif %s != 0 { q.Set(\"%s\", %s) }\n", varName, name, replacement))
				needStrconv = true
			} else if goType == "bool" {
				replacement := fmt.Sprintf("strconv.FormatBool(%s)", varName)
				queryBuild.WriteString(fmt.Sprintf("\tq.Set(\"%s\", %s)\n", name, replacement))
				needStrconv = true
			}
		} else if in == "header" {
			queryBuild.WriteString(fmt.Sprintf("\t// header param %s available in %s variable\n", name, varName))
		}
	}

	// build import block
	var importLines strings.Builder
	for imp := range imports {
		if imp == "strconv" && !needStrconv {
			continue
		}
		importLines.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
	}

	// construct the final command source code
	cmdCode := fmt.Sprintf(`package cmd

import (
%s
)
func New%sCmd() *cobra.Command {
%s
    cmd := &cobra.Command{
        Use:   "%s",
        Short: "%s",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := config.Load("%s")
            pathWithParams := "%s"
%s
            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
%s
            u.RawQuery = q.Encode()
            fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()
%s
            req, err := http.NewRequest("%s", fullUrl, bodyReader)
            if err != nil {
                return err
            }
			%s
			%s
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                return err
            }
            var pretty interface{}
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
%s
    return cmd
}
`, importLines.String(), cmdName, varDecls.String(), cmdName, op.Summary, moduleName, path, pathReplacements.String(), queryBuild.String(), bodyHandling, strings.ToUpper(method), headerHandling, authCode, flagsSetup.String())

	pathFile := filepath.Join(outputDir, "cmd", strings.ToLower(cmdName)+".go")
	return os.WriteFile(pathFile, []byte(cmdCode), 0644)
}

func sanitizeVar(s string) string {
	if s == "" {
		return "param"
	}
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	out := b.String()
	if out == "" {
		out = "param"
	}
	// if it starts with a digit, prefix with underscore
	if out[0] >= '0' && out[0] <= '9' {
		out = "_" + out
	}
	return out
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
