
package main

import (
	"petcli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewDelete_Store_Order_OrderIdCmd())
	rootCmd.AddCommand(cmd.NewGet_Store_Order_OrderIdCmd())
	rootCmd.AddCommand(cmd.NewGet_User_LogoutCmd())
	rootCmd.AddCommand(cmd.NewGet_User_UsernameCmd())
	rootCmd.AddCommand(cmd.NewPut_User_UsernameCmd())
	rootCmd.AddCommand(cmd.NewDelete_User_UsernameCmd())
	rootCmd.AddCommand(cmd.NewPost_PetCmd())
	rootCmd.AddCommand(cmd.NewPut_PetCmd())
	rootCmd.AddCommand(cmd.NewGet_Pet_FindByTagsCmd())
	rootCmd.AddCommand(cmd.NewDelete_Pet_PetIdCmd())
	rootCmd.AddCommand(cmd.NewGet_Pet_PetIdCmd())
	rootCmd.AddCommand(cmd.NewPost_Pet_PetIdCmd())
	rootCmd.AddCommand(cmd.NewPost_Pet_PetId_UploadImageCmd())
	rootCmd.AddCommand(cmd.NewPost_UserCmd())
	rootCmd.AddCommand(cmd.NewGet_User_LoginCmd())
	rootCmd.AddCommand(cmd.NewPost_User_CreateWithListCmd())
	rootCmd.AddCommand(cmd.NewGet_Store_InventoryCmd())
	rootCmd.AddCommand(cmd.NewPost_Store_OrderCmd())
	rootCmd.AddCommand(cmd.NewGet_Pet_FindByStatusCmd())

	cobra.CheckErr(rootCmd.Execute())
}
