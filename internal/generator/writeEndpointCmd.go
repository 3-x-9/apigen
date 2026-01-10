package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func writeEndpointCmd(outputDir string, moduleName string, cmdName string, op *openapi3.Operation, path, method string, authType string, authHeader string) error {
	// Build parameter-driven code pieces from the operation parameters
	var imports = map[string]bool{
		"fmt":                                true,
		"github.com/spf13/cobra":             true,
		fmt.Sprintf("%s/config", moduleName): true,
		"net/http":                           true,
		"io":                                 true,
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

	bodyHandling, headerHandling, debugHandling := buildBodyHeaderHandling(method, op, imports, &varDecls, &flagsSetup)

	for _, pRef := range op.Parameters {
		if pRef == nil || pRef.Value == nil {
			continue
		}
		p := pRef.Value
		name := p.Name
		in := p.In

		// init and default values to string
		goType := "string"
		flagFunc := "StringVar"

		varName := sanitizeVar(name)

		// assign correct goType and flagFunc based on parameter schema
		goType, flagFunc = FlagVars(goType, flagFunc, p, imports)

		// determine default value
		defaultVal := defaultForType(goType)

		varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))

		flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s %s parameter\")\n", flagFunc, varName, name, defaultVal, in, name))

		if pRef.Value.Required || p.In == "path" {
			flagsSetup.WriteString(fmt.Sprintf("\tcmd.MarkFlagRequired(\"%s\")\n", name))
		}

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
			} else if strings.HasPrefix(goType, "[]") {
				queryBuild.WriteString(fmt.Sprintf("\tfor _, v := range %s { q.Add(\"%s\", fmt.Sprintf(\"%%v\", v)) }\n", varName, name))
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
%s
			

            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                return err
            }
            var pretty interface{}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				fmt.Println("Request failed:")
				fmt.Printf("%%-15s: %%s\n", "Error", resp.Status)
				fmt.Printf("%%-15s: %%s\n", "URL", resp.Request.URL.String())
				fmt.Printf("%%-15s: %%s\n", "METHOD", resp.Request.Method)
				fmt.Println("----------------")
}
			if strings.Contains(resp.Header.Get("Content-Type"), "json") {
            if err := json.Unmarshal(body, &pretty); err != nil {
                return err
            }
            prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
            if err != nil {
                return err
            }
            fmt.Println("Response body:\n" + string(prettyJSON))
			} else {
			 	fmt.Println("Response body:\n" + string(body))
			}
            return nil
        },
    }
%s
    return cmd
}
`, importLines.String(), cmdName, varDecls.String(), cmdName, op.Summary, moduleName, path, pathReplacements.String(), queryBuild.String(), bodyHandling, strings.ToUpper(method), headerHandling, authCode, debugHandling, flagsSetup.String())

	pathFile := filepath.Join(outputDir, "cmd", strings.ToLower(cmdName)+".go")
	return os.WriteFile(pathFile, []byte(cmdCode), 0644)
}

func buildBodyHeaderHandling(method string, op *openapi3.Operation, imports map[string]bool, varDecls *strings.Builder, flagsSetup *strings.Builder) (string, string, string) {
	var bodyHandling string
	var headerHandling string
	var debugHandling string
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
		varDecls.WriteString("\tvar body string\n")
		flagsSetup.WriteString("\tcmd.Flags().StringVarP(&body, \"body\", \"b\", \"\", \"Request body (raw JSON, @filename, or '-' for stdin)\")\n")
		imports["bytes"] = true
		imports["os"] = true
		bodyHandling = `
			// prepare request body (allow raw JSON, @filename, or '-' for stdin)
			var bodyReader io.Reader
			if body != "" {
				if strings.HasPrefix(body, "@") {
					fname := strings.TrimPrefix(body, "@")
					var data []byte
					var err error
					if fname == "-" {
						data, err = io.ReadAll(os.Stdin)
						if err != nil {
							return err
						}
					} else {
						data, err = os.ReadFile(fname)
						if err != nil {
							return err
						}
					}
					bodyReader = bytes.NewReader(data)
				} else if body == "-" {
					data, err := io.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					bodyReader = bytes.NewReader(data)
				} else {
					bodyReader = strings.NewReader(body)
				}
			}
`
		debugHandling = buildDebugHandling(method)

		headerHandling = fmt.Sprintf(`
			if bodyReader != nil {
			    ct := contentType
				if ct == "" {
					ct = "%s"}
				req.Header.Set("Content-Type", ct)
			}
`, chooseRequestContentType(op))
		varDecls.WriteString("\tvar contentType string\n")
		flagsSetup.WriteString("\tcmd.Flags().StringVar(&contentType, \"content-type\", \"\", \"Content-Type header for the request body\")\n")
	} else if strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "DELETE" {
		// no body for GET/DELETE
		imports["io"] = true
		bodyHandling = `
			var bodyReader io.Reader = nil
`
		debugHandling = buildDebugHandling(method)

		headerHandling = ""
	}
	return bodyHandling, headerHandling, debugHandling
}

func buildDebugHandling(method string) string {
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
		return `				
			if Debug {
				fmt.Println("---DEBUG INFO---")
				fmt.Printf("%-15s: %s\n", "Request Method", req.Method)
				fmt.Printf("%-15s: %s\n", "URL", req.URL.String())
				fmt.Printf("%-15s: %v\n", "Headers", req.Header)
				if bodyReader != nil {
					data, err := io.ReadAll(bodyReader)
					if err != nil {
						return err
					}

					var parsed interface{}
					if json.Unmarshal(data, &parsed) == nil {
						prettyDebugJSON, err := json.MarshalIndent(parsed, "", "  ")
						if err != nil {
							return err
						}
					fmt.Printf("Request Body:\n%s\n", prettyDebugJSON)
					} else {
						fmt.Printf("Request Body:\n%s\n", string(data))
					}
					bodyReader = bytes.NewReader(data) // reset bodyReader
					
					} else {
						fmt.Printf("Request Body: (empty)\n")
					}
				fmt.Println("----------------")
			}`
	} else if strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "DELETE" {
		return `
			if Debug {
				fmt.Println("---DEBUG INFO---")
				fmt.Printf("%-15s: %s\n", "Request Method", req.Method)
				fmt.Printf("%-15s: %s\n", "URL", req.URL.String())
				fmt.Printf("%-15s: %v\n", "Headers", req.Header)
				fmt.Println("----------------")
			}`
	} else {
		return ""
	}
}

func defaultForType(goType string) string {
	switch goType {
	case "int":
		return "0"
	case "bool":
		return "false"
	case "[]int", "[]bool", "[]string":
		return "nil"
	default:
		return "\"\""
	}
}

func FlagVars(goType string, flagFunc string, p *openapi3.Parameter, imports map[string]bool) (string, string) {
	if p.Schema != nil && p.Schema.Value != nil && p.Schema.Value.Type != nil {
		switch (*p.Schema.Value.Type)[0] {
		case "array":
			if p.Schema.Value.Items != nil &&
				p.Schema.Value.Items.Value != nil &&
				p.Schema.Value.Items.Value.Type != nil &&
				len(*p.Schema.Value.Items.Value.Type) > 0 {
				switch (*p.Schema.Value.Items.Value.Type)[0] {
				case "integer":
					goType = "[]int"
					flagFunc = "IntSliceVar"
					imports["strconv"] = true
				case "boolean":
					goType = "[]bool"
					flagFunc = "BoolSliceVar"
					imports["strconv"] = true
				default:
					goType = "[]string"
					flagFunc = "StringSliceVar"
				}
			}

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
	return goType, flagFunc
}
