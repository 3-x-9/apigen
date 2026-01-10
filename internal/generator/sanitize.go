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

func sanitizeCommandName(path, method string) string {
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	parts := strings.Split(method+"_"+path, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		r := []rune(p)
		r[0] = unicode.ToUpper(r[0])
		parts[i] = string(r)
	}
	return strings.Join(parts, "_")
}
