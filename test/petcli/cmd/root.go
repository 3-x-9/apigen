
	package cmd

	import (
		"github.com/spf13/cobra"
	)

	func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cli",
		Short: "CLI is a command-line tool to interact with the API",
		}
		return cmd
		}
