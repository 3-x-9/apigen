package generator

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func buildBodyHeaderHandling(method string, op *openapi3.Operation, imports map[string]bool, varDecls *strings.Builder, flagsSetup *strings.Builder,
	schemaBodyBuild *strings.Builder, requiredChecks *strings.Builder) (string, string, string) {
	var bodyHandling string
	var headerHandling string
	var debugHandling string
	if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
		varDecls.WriteString("\tvar body string\n")
		flagsSetup.WriteString("\tcmd.Flags().StringVarP(&body, \"body\", \"b\", \"\", \"Request body (raw JSON, @filename, or '-' for stdin)\")\n")
		bodyHandling = fmt.Sprintf(`
			// prepare request body (allow raw JSON, @filename, or '-' for stdin)
			var bodyReader io.Reader
			var err error
			if body != "" {
				bodyReader, err = utils.GetBodyReader(body)
				if err != nil {
					return err
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
				utils.DebugPrintRequest(req, &bodyReader)
			}`
	} else if strings.ToUpper(method) == "GET" || strings.ToUpper(method) == "DELETE" {
		return `
			if Debug {
				utils.DebugPrintRequest(req, &bodyReader)
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
	SchemaBodyFlag *strings.Builder, validationBuild *strings.Builder, imports map[string]bool) {
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
			path := []string{}
			itterateProperties(schema, flagsSetup, varDecls, SchemaBodyFlag, path, validationBuild, imports)
		}
	}
	SchemaBodyFlag.WriteString(`if err := utils.CheckBody(body, bodyObj, &bodyReader); err != nil {
    return err
}`)
}

func buildRequestBodyFromSchemaFlags(schemaBodyBuild *strings.Builder, varName string, defaultVal string, path []string) []string {

	schemaBodyBuild.WriteString(fmt.Sprintf(`
	if %s != %s {`, varName, defaultVal))
	curr := "bodyObj"
	for i := 0; i < len(path)-1; i++ {
		key := path[i]
		schemaBodyBuild.WriteString(fmt.Sprintf(`if _, ok := %s["%s"]; !ok {
			%s["%s"] = map[string]interface{}{}
	}
		`, curr, key, curr, key))
		curr = fmt.Sprintf(`%s["%s"].(map[string]interface{})`, curr, key)
	}
	schemaBodyBuild.WriteString(fmt.Sprintf(`
		%s["%s"] = %s
		}
		`, curr, path[len(path)-1], varName))
	return path
}

func itterateProperties(schema *openapi3.Schema, flagsSetup *strings.Builder, varDecls *strings.Builder, SchemaBodyFlag *strings.Builder,
	path []string, validationBuild *strings.Builder, imports map[string]bool) {
	for propName, propRef := range schema.Properties {

		childPath := append(path, propName)

		if propRef != nil && propRef.Value != nil && propRef.Value.Type != nil && len(*propRef.Value.Type) > 0 {
			if (*propRef.Value.Type)[0] == "object" {

				itterateProperties(propRef.Value, flagsSetup, varDecls, SchemaBodyFlag, childPath, validationBuild, imports)
				continue
			}
		}

		propSchema := propRef.Value
		if propSchema == nil || propSchema.Type == nil {
			continue
		}

		varName := sanitizeVar(strings.Join(childPath, "_"))

		goType, flagFunc := mapSchemaToFlag(propSchema)

		varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))

		desc := propSchema.Description
		if desc == "" {
			desc = propName + " parameter"
		}

		if len(propSchema.Enum) > 0 {
			desc += fmt.Sprintf(" (one of: %v)", propSchema.Enum)
		}

		flagName := strings.Join(childPath, "-")
		flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s\")\n",
			flagFunc,
			varName,
			"body-"+flagName,
			defaultForType(goType),
			desc))
		_ = buildRequestBodyFromSchemaFlags(SchemaBodyFlag, varName, defaultForType(goType), childPath)

		// Check if this property is required
		for _, req := range schema.Required {
			if req == propName {
				flagsSetup.WriteString(fmt.Sprintf("\tcmd.MarkFlagRequired(\"%s\")\n", "body-"+flagName))
			}
		}

		generateEnumCheck(validationBuild, varName, propSchema.Enum, imports)
	}

}

func generateEnumCheck(checks *strings.Builder, varName string, enum []interface{}, imports map[string]bool) {
	if len(enum) == 0 {
		return
	}
	imports["fmt"] = true
	var enumValues []string
	for _, v := range enum {
		enumValues = append(enumValues, fmt.Sprintf("%v", v))
	}

	checks.WriteString(fmt.Sprintf(`
	if fmt.Sprintf("%%v", %s) != "" {
		valid := false
		allowed := []string{"%s"}
		for _, a := range allowed {
			if fmt.Sprintf("%%v", %s) == a {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid value for %s: %%v (allowed: %%v)", %s, allowed)
		}
	}
`, varName, strings.Join(enumValues, "\", \""), varName, varName, varName))
}
