package cmd

import (
	"github.com/spf13/cobra"
)

func NewStoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store",
		Short: "Commands related to store",
	}
}
