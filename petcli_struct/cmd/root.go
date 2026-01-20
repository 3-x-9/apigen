package cmd

import (
	"github.com/spf13/cobra"
)

var Debug bool
var Env string
var Output string

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "petcli_struct",
		Short: "petcli_struct is a command-line tool to interact with the API",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: false,
			HiddenDefaultCmd:  false,
		},
	}
	
	cmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug mode Show request/response details")
	cmd.PersistentFlags().StringVar(&Env, "env", "", "Environment to use (e.g. production, staging)")
	cmd.PersistentFlags().StringVar(&Output, "output", "pretty", "Output format (pretty, json, table, csv)")

	return cmd
}
