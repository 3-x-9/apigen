package cmd

import (
	"fmt"

	generator "github.com/3-x-9/apigen/internal/generator"
	"github.com/spf13/cobra"
)

var installCompletion bool

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate API code from OpenAPI specification",
	Long:  `Generate API code from OpenAPI specification using specified language and options`,
	RunE: func(cmd *cobra.Command, args []string) error {
		specPath, err := cmd.Flags().GetString("spec")
		if err != nil {
			return err
		}
		outputDir, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}
		moduleName, err := cmd.Flags().GetString("module")
		if err != nil {
			return err
		}

		installCompletion, err := cmd.Flags().GetBool("installCompletion")
		if installCompletion {
			if err := generator.InstallBashCompletion(cmd, moduleName); err != nil {
				fmt.Println(err)
				return err
			}
		}

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
	generateCmd.MarkFlagRequired("module")
	generateCmd.Flags().BoolVar(&installCompletion, "install-completion", false, "install a completion file to /etc/bash_completion.d/<module>")
}
