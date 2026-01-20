package cmd

import (
	"github.com/spf13/cobra"
)

func NewMiscCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "misc",
		Short: "Commands related to Misc",
	}
}
