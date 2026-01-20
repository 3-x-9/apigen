
package main

import (
	"petcli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	userTag := cmd.NewUserCmd()
	userTag.Use = "user"
	rootCmd.AddCommand(userTag)
	userTag.AddCommand(cmd.NewGetUserLoginCmd())
	userTag.AddCommand(cmd.NewGetUserUsernameCmd())
	userTag.AddCommand(cmd.NewPutUserUsernameCmd())
	userTag.AddCommand(cmd.NewDeleteUserUsernameCmd())
	userTag.AddCommand(cmd.NewPostUserCreateWithListCmd())
	userTag.AddCommand(cmd.NewGetUserLogoutCmd())
	userTag.AddCommand(cmd.NewPostUserCmd())
	petTag := cmd.NewPetCmd()
	petTag.Use = "pet"
	rootCmd.AddCommand(petTag)
	petTag.AddCommand(cmd.NewGetPetFindByStatusCmd())
	petTag.AddCommand(cmd.NewPostPetCmd())
	petTag.AddCommand(cmd.NewPutPetCmd())
	petTag.AddCommand(cmd.NewGetPetFindByTagsCmd())
	petTag.AddCommand(cmd.NewGetPetPetIdCmd())
	petTag.AddCommand(cmd.NewPostPetPetIdCmd())
	petTag.AddCommand(cmd.NewDeletePetPetIdCmd())
	petTag.AddCommand(cmd.NewPostPetPetIdUploadImageCmd())
	storeTag := cmd.NewStoreCmd()
	storeTag.Use = "store"
	rootCmd.AddCommand(storeTag)
	storeTag.AddCommand(cmd.NewGetStoreInventoryCmd())
	storeTag.AddCommand(cmd.NewDeleteStoreOrderOrderIdCmd())
	storeTag.AddCommand(cmd.NewGetStoreOrderOrderIdCmd())
	storeTag.AddCommand(cmd.NewPostStoreOrderCmd())


	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	cobra.CheckErr(rootCmd.Execute())
}
