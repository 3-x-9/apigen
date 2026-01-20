
package main

import (
	"testcli_list/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	miscTag := cmd.NewMiscCmd()
	miscTag.Use = "misc"
	rootCmd.AddCommand(miscTag)
	miscTag.AddCommand(cmd.NewGetItemsCmd())


	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	cobra.CheckErr(rootCmd.Execute())
}
