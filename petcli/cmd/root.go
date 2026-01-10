
	package cmd

	import (
		"github.com/spf13/cobra"
	)

	var Debug bool

	func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "petcli",
		Short: "petcli is a command-line tool to interact with the API",
		}
		cmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug mode Show request/response details")
		return cmd
		}
