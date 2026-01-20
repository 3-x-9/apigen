package generator

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type AuthScheme struct {
	Name       string // The key in securitySchemes
	Type       string // "http", "apiKey"
	Scheme     string // "bearer", "basic" etc (only for http)
	In         string // "header", "query", "cookie" (for apiKey)
	HeaderName string // Actual header name (for apiKey) or "Authorization" (for bearer)
}

func detectAuth(doc *openapi3.T) map[string]AuthScheme {
	schemes := make(map[string]AuthScheme)
	if doc == nil || doc.Components == nil || doc.Components.SecuritySchemes == nil {
		return schemes
	}
	for name, schemeRef := range doc.Components.SecuritySchemes {
		if schemeRef == nil || schemeRef.Value == nil {
			continue
		}
		s := schemeRef.Value

		switch s.Type {
		case "http":
			if s.Scheme == "bearer" {
				schemes[name] = AuthScheme{
					Name:       name,
					Type:       "http",
					Scheme:     "bearer",
					In:         "header",
					HeaderName: "Authorization",
				}
			}
		case "oauth2":
			schemes[name] = AuthScheme{
				Name:       name,
				Type:       "http",
				Scheme:     "bearer",
				In:         "header",
				HeaderName: "Authorization",
			}
		case "apiKey":
			schemes[name] = AuthScheme{
				Name:       name,
				Type:       "apiKey",
				In:         s.In,
				HeaderName: s.Name,
			}
		}
	}
	return schemes
}
