package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

func writeTagCmd(outputDir string, tag string) error {
	sanitizedTag := sanitizeTagName(tag)

	var buf bytes.Buffer
	data := struct {
		Tag         string
		Use         string
		TagOriginal string
	}{
		Tag:         sanitizedTag,
		Use:         sanitizeTagCLIName(tag),
		TagOriginal: tag,
	}

	if err := TagTmpl.Execute(&buf, data); err != nil {
		return err
	}

	path := filepath.Join(outputDir, "cmd", strings.ToLower(sanitizedTag)+"Tag.go")
	return os.WriteFile(path, buf.Bytes(), 0644)
}
