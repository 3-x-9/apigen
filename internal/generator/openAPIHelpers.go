package generator

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

func chooseRequestContentType(op *openapi3.Operation) string {
	if op == nil || op.RequestBody == nil || op.RequestBody.Value == nil {
		return ""
	}
	prefs := []string{
		"application/json",
		"+json", // match suffix
		"application/xml",
		"application/x-www-form-urlencoded",
		"multipart/form-data",
	}
	// direct matches first
	for _, p := range prefs {
		for ct := range op.RequestBody.Value.Content {
			if p == "+json" {
				if strings.HasSuffix(ct, "+json") {
					return ct
				}
			} else if ct == p {
				return ct
			}
		}
	}
	// fallback: return the first available content type
	for ct := range op.RequestBody.Value.Content {
		return ct
	}
	return ""
}
