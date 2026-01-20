package generator

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func InstallBashCompletion(rootCmd *cobra.Command, moduleName string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dirPath := filepath.Join(home + "./bash_completion.d")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	path := filepath.Join(dirPath, moduleName)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return rootCmd.GenBashCompletion(f)
}
