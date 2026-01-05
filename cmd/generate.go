package cmd

import (
	"github.com/3-x-9/apigen/internal/generator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate API code from OpenAPI specification",
	Long:  `Generate API code from OpenAPI specification using specified language and options`,
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath, _ := cmd.Flags().GetString("spec")
		outputDir, _ := cmd.Flags().GetString("out")
		moduleName, _ := cmd.Flags().GetString("module")

		gen := generator.NewGenerator()
		return gen.Generate(specPath, outputDir, moduleName)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("spec", "s", "", "Path to the OpenAPI specification file (required)")
	generateCmd.MarkFlagRequired("spec")
	generateCmd.Flags().StringP("out", "o", "./cli-out", "Output directory for the generated code")
	generateCmd.Flags().StringP("module", "m", "", "Module name for the generated code")
}
