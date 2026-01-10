package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeRootCmd(outputDir string, moduleName string) error {
	rootCode := fmt.Sprintf(`
	package cmd

	import (
		"github.com/spf13/cobra"
	)

	var Debug bool

	func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "%s",
		Short: "%s is a command-line tool to interact with the API",
		}
		cmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug mode Show request/response details")
		return cmd
		}
`, moduleName, moduleName)
	path := filepath.Join(outputDir, "cmd", "root.go")
	return os.WriteFile(path, []byte(rootCode), 0644)
}
