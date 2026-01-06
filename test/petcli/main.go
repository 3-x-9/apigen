
package main

import (
	"petcli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewGet_user_usernameCmd())
	rootCmd.AddCommand(cmd.NewPut_user_usernameCmd())
	rootCmd.AddCommand(cmd.NewDelete_user_usernameCmd())
	rootCmd.AddCommand(cmd.NewPost_petCmd())
	rootCmd.AddCommand(cmd.NewPut_petCmd())
	rootCmd.AddCommand(cmd.NewPost_pet_petId_uploadImageCmd())
	rootCmd.AddCommand(cmd.NewGet_store_order_orderIdCmd())
	rootCmd.AddCommand(cmd.NewDelete_store_order_orderIdCmd())
	rootCmd.AddCommand(cmd.NewGet_pet_findByStatusCmd())
	rootCmd.AddCommand(cmd.NewGet_user_loginCmd())
	rootCmd.AddCommand(cmd.NewGet_user_logoutCmd())
	rootCmd.AddCommand(cmd.NewPost_pet_petIdCmd())
	rootCmd.AddCommand(cmd.NewDelete_pet_petIdCmd())
	rootCmd.AddCommand(cmd.NewGet_pet_petIdCmd())
	rootCmd.AddCommand(cmd.NewPost_user_createWithListCmd())
	rootCmd.AddCommand(cmd.NewGet_pet_findByTagsCmd())
	rootCmd.AddCommand(cmd.NewGet_store_inventoryCmd())
	rootCmd.AddCommand(cmd.NewPost_store_orderCmd())
	rootCmd.AddCommand(cmd.NewPost_userCmd())

	cobra.CheckErr(rootCmd.Execute())
}
