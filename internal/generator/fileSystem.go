package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func createFolders(outputDir string) error {
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "cmd"),
		filepath.Join(outputDir, "config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

	}
	return nil
}

func writeGoMod(outputDir, moduleName string) error {
	goModContent := fmt.Sprintf(`module %s

go 1.22

require (
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.19.0
)
`, moduleName)

	path := filepath.Join(outputDir, "go.mod")
	return os.WriteFile(path, []byte(goModContent), 0644)
}

func writeMain(outputDir string, moduleName string, cmds []string) error {
	cmdsInit := ""
	for _, c := range cmds {
		cmdsInit += fmt.Sprintf("\trootCmd.AddCommand(cmd.New%sCmd())\n", c)
	}

	mainCode := fmt.Sprintf(`
package main

import (
	"%s/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
%s
	cobra.CheckErr(rootCmd.Execute())
}
`, moduleName, cmdsInit)

	path := filepath.Join(outputDir, "main.go")
	return os.WriteFile(path, []byte(mainCode), 0644)
}
