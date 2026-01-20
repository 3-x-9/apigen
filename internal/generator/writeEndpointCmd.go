package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func writeEndpointCmd(outputDir string, moduleName string, goName string, cliName string, op *openapi3.Operation, path, method string,
	schemes map[string]AuthScheme, globalSecurity *openapi3.SecurityRequirements) error {
	// Build parameter-driven code pieces from the operation parameters
	var imports = map[string]bool{
		"github.com/spf13/cobra":             true,
		fmt.Sprintf("%s/config", moduleName): true,
		fmt.Sprintf("%s/utils", moduleName):  true,
		"net/http":                           true,
		"io":                                 true,
		"net/url":                            true,
		"strings":                            true,
	}

	respModel, isArray := detectResponseModel(op)
	if respModel != "" {
		imports[fmt.Sprintf("%s/models", moduleName)] = true
	}

	var varDecls strings.Builder
	var flagsSetup strings.Builder
	var pathReplacements strings.Builder
	var queryBuild strings.Builder
	var schemaBodyBuild strings.Builder
	var validationBuild strings.Builder

	var headerBuild strings.Builder
	var cookieBuild strings.Builder

	// Determine security requirements: Op level overrides Global
	security := globalSecurity
	if op.Security != nil {
		security = op.Security
	}

	// Find the first matching scheme

	authCode := findScheme(security, schemes, &cookieBuild, &queryBuild)

	generateBodyFlagsFromSchema(op, &flagsSetup, &varDecls, &schemaBodyBuild, &validationBuild, imports)

	bodyHandling, headerHandling, _ := buildBodyHeaderHandling(method, op, imports, &varDecls, &flagsSetup, &schemaBodyBuild, &validationBuild)

	buildPathParams(*op, &varDecls, &flagsSetup, &pathReplacements, &queryBuild, &headerBuild, &cookieBuild, &validationBuild, imports)

	var importList []string
	for imp := range imports {
		importList = append(importList, imp)
	}

	writePkgUtil(outputDir)

	err := buildCmdCode(CmdConfig{
		Method:           strings.ToUpper(method),
		GoName:           goName,
		CommandName:      cliName,
		ModuleName:       moduleName,
		Path:             path,
		OutputDir:        outputDir,
		Short:            op.Summary,
		VarDecls:         varDecls.String(),
		PathReplacements: pathReplacements.String(),
		QueryBuild:       queryBuild.String(),
		HeaderBuild:      headerBuild.String(),
		CookieBuild:      cookieBuild.String(),
		FlagsSetup:       flagsSetup.String(),
		Imports:          importList,
		BodyHandling:     bodyHandling,
		HeaderHandling:   headerHandling,
		AuthCode:         authCode,
		Validation:       validationBuild.String(),
		ResponseModel:    respModel,
		IsArray:          isArray,
	})
	return err
}

func detectResponseModel(op *openapi3.Operation) (string, bool) {
	if op.Responses == nil {
		return "", false
	}
	// Check for 200, 201, or default success
	successCodes := []string{"200", "201", "202", "204", "default"}
	for _, code := range successCodes {
		respRef := op.Responses.Value(code)
		if respRef == nil || respRef.Value == nil || respRef.Value.Content == nil {
			continue
		}
		content := respRef.Value.Content.Get("application/json")
		if content == nil || content.Schema == nil {
			continue
		}
		schemaRef := content.Schema
		if schemaRef.Ref != "" {
			parts := strings.Split(schemaRef.Ref, "/")
			return toGoName(parts[len(parts)-1]), false
		}
		// If it's an array of references
		if schemaRef.Value != nil && schemaRef.Value.Type != nil && len(*schemaRef.Value.Type) > 0 && (*schemaRef.Value.Type)[0] == "array" {
			if schemaRef.Value.Items != nil && schemaRef.Value.Items.Ref != "" {
				parts := strings.Split(schemaRef.Value.Items.Ref, "/")
				return toGoName(parts[len(parts)-1]), true
			}
		}
	}
	return "", false
}

type CmdConfig struct {
	Method           string
	GoName           string
	CommandName      string
	ModuleName       string
	Path             string
	OutputDir        string
	Short            string
	VarDecls         string
	PathReplacements string
	QueryBuild       string
	HeaderBuild      string
	CookieBuild      string
	FlagsSetup       string
	Imports          []string
	BodyHandling     string
	HeaderHandling   string
	AuthCode         string
	Validation       string
	ResponseModel    string
	IsArray          bool
}

func buildCmdCode(cfg CmdConfig) error {
	var buf bytes.Buffer
	if err := EndpointTmpl.Execute(&buf, cfg); err != nil {
		return err
	}

	pathFile := filepath.Join(cfg.OutputDir, "cmd", strings.ToLower(cfg.GoName)+".go")
	return os.WriteFile(pathFile, buf.Bytes(), 0644)
}

func buildPathParams(op openapi3.Operation, varDecls *strings.Builder, flagsSetup *strings.Builder, pathReplacements *strings.Builder, queryBuild *strings.Builder,
	headerBuild *strings.Builder, cookieBuild *strings.Builder, validationBuild *strings.Builder, imports map[string]bool) {
	for _, pRef := range op.Parameters {
		if pRef == nil || pRef.Value == nil {
			continue
		}
		p := pRef.Value
		name := p.Name
		in := p.In

		// detect explode in array
		explode := false
		if p.Style != "" {
			if p.Style == "form" && p.Explode != nil {
				explode = *p.Explode
			}
		}

		separator := ","
		switch p.Style {
		case "pipeDelimited":
			separator = "|"
		case "spaceDelimited":
			separator = " "
		}

		// init and default values to string
		goType := "string"
		flagFunc := "StringVar"

		varName := sanitizeVar(name)

		// assign correct goType and flagFunc based on parameter schema
		goType, flagFunc = FlagVars(goType, flagFunc, p, imports)

		// determine default value
		defaultVal := defaultForType(goType)

		switch in {
		case "path":
			varName += "Path"
			varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))
			flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s %s parameter\")\n", flagFunc, varName, name, defaultVal, in, name))
		case "query":
			varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))
			flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s %s parameter\")\n", flagFunc, varName, name, defaultVal, in, name))
		case "header":
			varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))
			flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s %s parameter\")\n", flagFunc, varName, name, defaultVal, in, name))
		case "cookie":
			varDecls.WriteString(fmt.Sprintf("\tvar %s %s\n", varName, goType))
			flagsSetup.WriteString(fmt.Sprintf("\tcmd.Flags().%s(&%s, \"%s\", %s, \"%s %s parameter\")\n", flagFunc, varName, name, defaultVal, in, name))
		default:
			continue
		}

		if pRef.Value.Required || p.In == "path" {
			flagsSetup.WriteString(fmt.Sprintf("\tcmd.MarkFlagRequired(\"%s\")\n", name))
		}
		if p.Schema != nil && p.Schema.Value != nil {
			generateEnumCheck(validationBuild, varName, p.Schema.Value.Enum, imports)
		}
		switch in {
		case "path":
			// replace path placeholders with param values
			replacement := varName
			switch goType {
			case "int":
				replacement = fmt.Sprintf("strconv.Itoa(%s)", varName)
				imports["strconv"] = true
			case "bool":
				replacement = fmt.Sprintf("strconv.FormatBool(%s)", varName)
				imports["strconv"] = true
			}
			pathReplacements.WriteString(fmt.Sprintf("\tpathWithParams = strings.ReplaceAll(pathWithParams, \"{%s}\", %s)\n", name, replacement))
		case "query":
			if goType == "string" {
				queryBuild.WriteString(fmt.Sprintf("\tif %s != \"\" { q.Set(\"%s\", %s) }\n", varName, name, varName))
			} else if goType == "int" {
				replacement := fmt.Sprintf("strconv.Itoa(%s)", varName)
				queryBuild.WriteString(fmt.Sprintf("\tif %s != 0 { q.Set(\"%s\", %s) }\n", varName, name, replacement))
				imports["strconv"] = true
			} else if goType == "bool" {
				replacement := fmt.Sprintf("strconv.FormatBool(%s)", varName)
				queryBuild.WriteString(fmt.Sprintf("\tq.Set(\"%s\", %s)\n", name, replacement))
				imports["strconv"] = true
			} else if strings.HasPrefix(goType, "[]") {
				imports["fmt"] = true
				if explode {
					queryBuild.WriteString(fmt.Sprintf("\tfor _, v := range %s { q.Add(\"%s\", fmt.Sprintf(\"%%v\", v)) }\n", varName, name))
				} else {
					queryBuild.WriteString(fmt.Sprintf("\tif %s != nil { q.Set(\"%s\", strings.Join(func() []string { res := []string{}; for _, v := range %s { res = append(res, fmt.Sprintf(\"%%v\", v)) }; return res }(), \"%s\")) }\n", varName, name, varName, separator))
				}
			}
		case "header":
			switch goType {
			case "string":
				headerBuild.WriteString(fmt.Sprintf("\tif %s != \"\" { req.Header.Set(\"%s\", %s) }\n", varName, name, varName))
			case "int":
				replacement := fmt.Sprintf("strconv.Itoa(%s)", varName)
				headerBuild.WriteString(fmt.Sprintf("\tif %s != 0 { req.Header.Set(\"%s\", %s) }\n", varName, name, replacement))
				imports["strconv"] = true
			case "bool":
				replacement := fmt.Sprintf("strconv.FormatBool(%s)", varName)
				headerBuild.WriteString(fmt.Sprintf("\treq.Header.Set(\"%s\", %s)\n", name, replacement))
				imports["strconv"] = true
			}
		case "cookie":
			switch goType {
			case "string":
				cookieBuild.WriteString(fmt.Sprintf("\tif %s != \"\" { req.AddCookie(&http.Cookie{Name: \"%s\", Value: %s}) }\n", varName, name, varName))
			case "int":
				replacement := fmt.Sprintf("strconv.Itoa(%s)", varName)
				cookieBuild.WriteString(fmt.Sprintf("\tif %s != 0 { req.AddCookie(&http.Cookie{Name: \"%s\", Value: %s}) }\n", varName, name, replacement))
				imports["strconv"] = true
			case "bool":
				replacement := fmt.Sprintf("strconv.FormatBool(%s)", varName)
				cookieBuild.WriteString(fmt.Sprintf("\treq.AddCookie(&http.Cookie{Name: \"%s\", Value: %s})\n", name, replacement))
				imports["strconv"] = true
			}
		}
	}
}

func findScheme(security *openapi3.SecurityRequirements, schemes map[string]AuthScheme, cookieBuild *strings.Builder,
	queryBuild *strings.Builder) string {
	var selectedScheme *AuthScheme
	var authCode string
	if security != nil {
		found := false
		for _, req := range *security {
			for name := range req {
				if s, ok := schemes[name]; ok {
					selectedScheme = &s
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if selectedScheme != nil {
		// e.g. "PetstoreAuth"
		fieldName := selectedScheme.Name
		if len(fieldName) > 0 {
			fieldName = strings.ToUpper(fieldName[:1]) + fieldName[1:]
		}
		fieldName += "Auth"

		if selectedScheme.Type == "http" && selectedScheme.Scheme == "bearer" {
			authCode = fmt.Sprintf(`
		if cfg.%s != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.%s)
		}
		`, fieldName, fieldName)
		} else if selectedScheme.Type == "apiKey" {
			switch selectedScheme.In {
			case "header":
				authCode = fmt.Sprintf(`
		if cfg.%s != "" {
			req.Header.Set("%s", cfg.%s)}`, fieldName, selectedScheme.HeaderName, fieldName)
			case "query":
				queryBuild.WriteString(fmt.Sprintf(`
		if cfg.%s != "" {
			q.Set("%s", cfg.%s)
		}`, fieldName, selectedScheme.HeaderName, fieldName))
			case "cookie":
				cookieBuild.WriteString(fmt.Sprintf(`
		if cfg.%s != "" {
			req.AddCookie(&http.Cookie{Name: "%s", Value: cfg.%s})
		}
`, fieldName, selectedScheme.HeaderName, fieldName))
			}
		}
	}
	return authCode
}
