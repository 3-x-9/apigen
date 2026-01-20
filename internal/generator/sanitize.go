package generator

import (
	"strings"
	"unicode"
)

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

func toGoName(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for i, p := range parts {
		r := []rune(p)
		r[0] = unicode.ToUpper(r[0])
		parts[i] = string(r)
	}
	return strings.Join(parts, "")
}

func toCLIName(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for i, p := range parts {
		parts[i] = strings.ToLower(p)
	}
	return strings.Join(parts, "-")
}

func sanitizeCommandName(path, method string) string {
	cleanPath := strings.ReplaceAll(path, "{", "")
	cleanPath = strings.ReplaceAll(cleanPath, "}", "")
	return toGoName(method + "_" + cleanPath)
}

func sanitizeCLIName(path, method string) string {
	cleanPath := strings.ReplaceAll(path, "{", "")
	cleanPath = strings.ReplaceAll(cleanPath, "}", "")
	return toCLIName(method + "_" + cleanPath)
}

func sanitizeTagName(tag string) string {
	return toGoName(tag)
}

func sanitizeTagCLIName(tag string) string {
	return toCLIName(tag)
}
