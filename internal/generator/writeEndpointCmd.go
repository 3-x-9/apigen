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
	var schemaBodyBuild strings.Builder
	var needStrconv bool
	var authCode string
	var requiredChecks strings.Builder

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

	generateBodyFlagsFromSchema(op, &flagsSetup, &varDecls, &schemaBodyBuild, &requiredChecks)

	bodyHandling, headerHandling, debugHandling := buildBodyHeaderHandling(method, op, imports, &varDecls, &flagsSetup, &schemaBodyBuild, &requiredChecks)

	for _, pRef := range op.Parameters {
		if pRef == nil || pRef.Value == nil {
			continue
		}
		p := pRef.Value
		name := p.Name
		in := p.In

		// detect explode in arrayt
		explode := false
		if p.Style != "" {
			if p.Style == "form" && p.Explode != nil {
				explode = *p.Explode
			}
		}

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
				if explode {
					queryBuild.WriteString(fmt.Sprintf("\tfor _, v := range %s { q.Add(\"%s\", fmt.Sprintf(\"%%v\", v)) }\n", varName, name))
				} else {
					queryBuild.WriteString(fmt.Sprintf("\tif %s != nil { q.Set(\"%s\", strings.Join(func() []string { res := []string{}; for _, v := range %s { res = append(res, fmt.Sprintf(\"%%v\", v)) }; return res }(), \",\")) }\n", varName, name, varName))
				}
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

func buildBodyHeaderHandling(method string, op *openapi3.Operation, imports map[string]bool, varDecls *strings.Builder, flagsSetup *strings.Builder,
	schemaBodyBuild *strings.Builder, requiredChecks *strings.Builder) (string, string, string) {
	var bodyHandling string
	var headerHandling string
	var debugHandling string
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
		varDecls.WriteString("\tvar body string\n")
		flagsSetup.WriteString("\tcmd.Flags().StringVarP(&body, \"body\", \"b\", \"\", \"Request body (raw JSON, @filename, or '-' for stdin)\")\n")
		imports["bytes"] = true
		imports["os"] = true
		bodyHandling = fmt.Sprintf(`
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
					bodyReader = bytes.NewReader([]byte(body))
				}
			} else {
			 %s
			if body == "" {
				%s
				}	
			}
`, schemaBodyBuild, requiredChecks)
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

func mapSchemaToFlag(prop *openapi3.Schema) (string, string) {
	// Default to string
	goType := "string"
	flagFunc := "StringVar"
	if prop == nil || prop.Type == nil || len(*prop.Type) == 0 {
		return goType, flagFunc
	}
	switch (*prop.Type)[0] {
	case "array":
		if prop.Items != nil &&
			prop.Items.Value != nil &&
			prop.Items.Value.Type != nil &&
			len(*prop.Items.Value.Type) > 0 {

			switch (*prop.Items.Value.Type)[0] {
			case "integer":
				goType = "[]int"
				flagFunc = "IntSliceVar"
			case "boolean":
				goType = "[]bool"
				flagFunc = "BoolSliceVar"
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
	case "boolean":
		goType = "bool"
		flagFunc = "BoolVar"
	case "object":
		goType = "map[string]interface{}"
		flagFunc = "StringVar"
	default:
		// default to string for unknown types
		goType = "string"
		flagFunc = "StringVar"
	}
	return goType, flagFunc
}

func generateBodyFlagsFromSchema(op *openapi3.Operation, flagsSetup *strings.Builder, varDecls *strings.Builder,
	SchemaBodyFlag *strings.Builder, requiredChecks *strings.Builder) {
	if op.RequestBody == nil || op.RequestBody.Value == nil || op.RequestBody.Value.Content == nil {
		return
	}
	content := op.RequestBody.Value.Content.Get("application/json")
	if content == nil || content.Schema == nil || content.Schema.Value == nil {
		return
	}
	schema := content.Schema.Value
	if schema.Type != nil && len(*schema.Type) > 0 {
		SchemaBodyFlag.WriteString(`bodyObj := map[string]interface{}{}
		`)
		switch (*schema.Type)[0] {
		case "object":
			flagStr := "body"
			varStr := "body"
			itterateProperties(schema, flagsSetup, varDecls, SchemaBodyFlag, flagStr, varStr, requiredChecks)
		}
	}
	SchemaBodyFlag.WriteString(`if body == "" && len(bodyObj) <= 0 {
									return fmt.Errorf("request body is required (use either --body or body flags!!!)")
								}
								data, err := json.Marshal(bodyObj)
								if err != nil {
									return err
								}
								bodyReader = bytes.NewReader(data)
								`)
}

func buildRequestBodyFromSchemaFlags(schemaBodyBuild *strings.Builder, varName string, varType string) {
	schemaBodyBuild.WriteString(fmt.Sprintf(`
	if %s != %s {
		bodyObj["%s"] = %s
	}
`, varName, varType, varName, varName))
}

func isObjectSchema(s *openapi3.Schema) bool {
	if s == nil || s.Type == nil {
		return false
	}
	for _, t := range *s.Type {
		if t == "object" {
			return true
		}
	}
	return false
}

func itterateProperties(schema *openapi3.Schema, flagsSetup *strings.Builder, varDecls *strings.Builder, SchemaBodyFlag *strings.Builder,
	flagStr string, varStr string, requiredChecks *strings.Builder) {
	for propName, propRef := range schema.Properties {
		if propRef != nil && propRef.Value != nil && propRef.Value.Type != nil && len(*propRef.Value.Type) > 0 {

			if (*propRef.Value.Type)[0] == "object" {
				childFlagStr := fmt.Sprintf("%s-%s", flagStr, propName)
				childVarStr := fmt.Sprintf("%s_%s", varStr, propName)
				itterateProperties(propRef.Value, flagsSetup, varDecls, SchemaBodyFlag, childFlagStr, childVarStr, requiredChecks)
				continue
			}

		}
		propSchema := propRef.Value
		if propSchema == nil || propSchema.Type == nil {
			continue
		}

		varName := sanitizeVar(propName + "_" + varStr)

		goType, flagFunc := mapSchemaToFlag(propSchema)

		varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))

		desc := propSchema.Description
		if desc == "" {
			desc = propName + " parameter"
		}

		if len(propSchema.Enum) > 0 {
			desc += fmt.Sprintf(" (one of: %v)", propSchema.Enum)
		}

		propName = flagStr + "-" + propName
		flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s\")\n",
			flagFunc,
			varName,
			propName,
			defaultForType(goType),
			desc))

		buildRequestBodyFromSchemaFlags(SchemaBodyFlag, varName, defaultForType(goType))
	}
	generateRequiredChecks(schema, "bodyObj", requiredChecks)
}

func generateRequiredChecks(schema *openapi3.Schema, parentVar string, builder *strings.Builder) {
	for _, reqName := range schema.Required {
		propRef := schema.Properties[reqName]
		if propRef == nil || propRef.Value == nil {
			continue
		}
		reqName += "_body"
		childVar := fmt.Sprintf("%s[\"%s\"]", parentVar, reqName)
		if isObjectSchema(propRef.Value) {
			builder.WriteString(fmt.Sprintf(`
nested, ok := %s.(map[string]interface{})
if !ok {
    return fmt.Errorf("missing required object: %s")
}
`, childVar, reqName))
			generateRequiredChecks(propRef.Value, "nested", builder)
		} else {
			builder.WriteString(fmt.Sprintf(`
if _, ok := %s; !ok {
    return fmt.Errorf("missing required field: %s")
}
`, childVar, reqName))
		}
	}
}
