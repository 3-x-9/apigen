package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "apigen",
	Short: "API Generator is a tool to generate a code for an API CLI from OpenAPI specifications",
	Long:  `API Generator is a command-line tool that generates CLI code for an API CLI from OpenAPI specifications`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
