package cmd

import (
	"github.com/spf13/cobra"
)

func NewPetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pet",
		Short: "Commands related to pet",
	}
}
