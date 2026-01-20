
package main

import (
	"petcli_typed/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	petTag := cmd.NewPetCmd()
	petTag.Use = "pet"
	rootCmd.AddCommand(petTag)
	petTag.AddCommand(cmd.NewGetPetFindByTagsCmd())
	petTag.AddCommand(cmd.NewGetPetPetIdCmd())
	petTag.AddCommand(cmd.NewPostPetPetIdCmd())
	petTag.AddCommand(cmd.NewDeletePetPetIdCmd())
	petTag.AddCommand(cmd.NewPostPetCmd())
	petTag.AddCommand(cmd.NewPutPetCmd())
	petTag.AddCommand(cmd.NewPostPetPetIdUploadImageCmd())
	petTag.AddCommand(cmd.NewGetPetFindByStatusCmd())
	storeTag := cmd.NewStoreCmd()
	storeTag.Use = "store"
	rootCmd.AddCommand(storeTag)
	storeTag.AddCommand(cmd.NewGetStoreOrderOrderIdCmd())
	storeTag.AddCommand(cmd.NewDeleteStoreOrderOrderIdCmd())
	storeTag.AddCommand(cmd.NewGetStoreInventoryCmd())
	storeTag.AddCommand(cmd.NewPostStoreOrderCmd())
	userTag := cmd.NewUserCmd()
	userTag.Use = "user"
	rootCmd.AddCommand(userTag)
	userTag.AddCommand(cmd.NewPostUserCmd())
	userTag.AddCommand(cmd.NewGetUserLoginCmd())
	userTag.AddCommand(cmd.NewGetUserLogoutCmd())
	userTag.AddCommand(cmd.NewPostUserCreateWithListCmd())
	userTag.AddCommand(cmd.NewDeleteUserUsernameCmd())
	userTag.AddCommand(cmd.NewGetUserUsernameCmd())
	userTag.AddCommand(cmd.NewPutUserUsernameCmd())


	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	cobra.CheckErr(rootCmd.Execute())
}
