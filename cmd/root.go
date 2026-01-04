package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "apigen",
	Short: "API Generator is a tool to generate API code from OpenAPI specifications",
	Long:  `API Generator is a command-line tool that generates API code from OpenAP specifications`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
