package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func writeModels(outputDir string, doc *openapi3.T) error {
	var b strings.Builder
	b.WriteString("package models\n\n")

	for name, schemaRef := range doc.Components.Schemas {
		if schemaRef.Value == nil {
			continue
		}
		structCode := generateStruct(name, schemaRef.Value)
		b.WriteString(structCode + "\n")
	}

	modelsDir := filepath.Join(outputDir, "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(modelsDir, "models.go"), []byte(b.String()), 0644)
}

func generateStruct(name string, schema *openapi3.Schema) string {
	var b strings.Builder
	typeName := toGoName(name)
	b.WriteString(fmt.Sprintf("type %s struct {\n", typeName))

	for propName, propRef := range schema.Properties {
		if propRef.Value == nil {
			continue
		}
		fieldName := toGoName(propName)
		goType := mapSchemaToGoType(propRef)
		b.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, goType, propName))
	}

	b.WriteString("}\n")
	return b.String()
}

func mapSchemaToGoType(ref *openapi3.SchemaRef) string {
	if ref == nil {
		return "interface{}"
	}

	if ref.Ref != "" {
		parts := strings.Split(ref.Ref, "/")
		return toGoName(parts[len(parts)-1])
	}

	schema := ref.Value
	if schema == nil || schema.Type == nil || len(*schema.Type) == 0 {
		return "interface{}"
	}

	switch (*schema.Type)[0] {
	case "string":
		return "string"
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil {
			return "[]" + mapSchemaToGoType(schema.Items)
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}
