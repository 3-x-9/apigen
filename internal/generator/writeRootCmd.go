package generator

import (
	"bytes"
	"os"
	"path/filepath"
)

func writeRootCmd(outputDir string, moduleName string) error {
	var buf bytes.Buffer
	data := struct {
		ModuleName string
	}{
		ModuleName: moduleName,
	}

	if err := RootTmpl.Execute(&buf, data); err != nil {
		return err
	}
	path := filepath.Join(outputDir, "cmd", "root.go")
	return os.WriteFile(path, buf.Bytes(), 0644)
}
