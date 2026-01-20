package cmd

import (
	"github.com/spf13/cobra"
)

func NewUserCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "user",
		Short: "Commands related to user",
	}
}
