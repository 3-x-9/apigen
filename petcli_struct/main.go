
package main

import (
	"petcli_struct/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	petTag := cmd.NewPetCmd()
	petTag.Use = "pet"
	rootCmd.AddCommand(petTag)
	petTag.AddCommand(cmd.NewPostPetCmd())
	petTag.AddCommand(cmd.NewPutPetCmd())
	petTag.AddCommand(cmd.NewGetPetFindByStatusCmd())
	petTag.AddCommand(cmd.NewPostPetPetIdUploadImageCmd())
	petTag.AddCommand(cmd.NewGetPetPetIdCmd())
	petTag.AddCommand(cmd.NewPostPetPetIdCmd())
	petTag.AddCommand(cmd.NewDeletePetPetIdCmd())
	petTag.AddCommand(cmd.NewGetPetFindByTagsCmd())
	userTag := cmd.NewUserCmd()
	userTag.Use = "user"
	rootCmd.AddCommand(userTag)
	userTag.AddCommand(cmd.NewPostUserCmd())
	userTag.AddCommand(cmd.NewPostUserCreateWithListCmd())
	userTag.AddCommand(cmd.NewGetUserLoginCmd())
	userTag.AddCommand(cmd.NewGetUserLogoutCmd())
	userTag.AddCommand(cmd.NewGetUserUsernameCmd())
	userTag.AddCommand(cmd.NewPutUserUsernameCmd())
	userTag.AddCommand(cmd.NewDeleteUserUsernameCmd())
	storeTag := cmd.NewStoreCmd()
	storeTag.Use = "store"
	rootCmd.AddCommand(storeTag)
	storeTag.AddCommand(cmd.NewGetStoreInventoryCmd())
	storeTag.AddCommand(cmd.NewPostStoreOrderCmd())
	storeTag.AddCommand(cmd.NewDeleteStoreOrderOrderIdCmd())
	storeTag.AddCommand(cmd.NewGetStoreOrderOrderIdCmd())


	rootCmd.CompletionOptions.DisableDefaultCmd = false
	rootCmd.CompletionOptions.HiddenDefaultCmd = false

	cobra.CheckErr(rootCmd.Execute())
}
