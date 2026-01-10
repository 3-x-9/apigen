package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

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
