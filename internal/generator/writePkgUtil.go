package generator

import (
	"bytes"
	"os"
	"path/filepath"
)

func writePkgUtil(OutputDir string) error {
	var buf bytes.Buffer
	if err := UtilTmpl.Execute(&buf, nil); err != nil {
		return err
	}

	pathFile := filepath.Join(OutputDir, "utils", "utils.go")
	return os.WriteFile(pathFile, buf.Bytes(), 0644)
}
