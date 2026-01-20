package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func createFolders(outputDir string) error {
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "cmd"),
		filepath.Join(outputDir, "config"),
		filepath.Join(outputDir, "utils"),
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

func writeMain(outputDir string, moduleName string, tagToCmds map[string][]CommandInfo) error {
	cmdsInit := ""
	for tag, cmds := range tagToCmds {
		sanitizedTagGo := sanitizeTagName(tag)
		sanitizedTagCLI := sanitizeTagCLIName(tag)
		tagVar := strings.ToLower(sanitizedTagGo) + "Tag"
		cmdsInit += fmt.Sprintf("\t%s := cmd.New%sCmd()\n", tagVar, sanitizedTagGo)
		cmdsInit += fmt.Sprintf("\t%s.Use = \"%s\"\n", tagVar, sanitizedTagCLI)
		cmdsInit += fmt.Sprintf("\trootCmd.AddCommand(%s)\n", tagVar)
		for _, c := range cmds {
			cmdsInit += fmt.Sprintf("\t%s.AddCommand(cmd.New%sCmd())\n", tagVar, c.GoName)
		}
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

	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	cobra.CheckErr(rootCmd.Execute())
}
`, moduleName, cmdsInit)

	path := filepath.Join(outputDir, "main.go")
	return os.WriteFile(path, []byte(mainCode), 0644)
}
