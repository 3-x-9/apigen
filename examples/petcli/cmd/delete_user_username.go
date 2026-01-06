package cmd

import (
	"petcli/config"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"

)
func NewDelete_user_usernameCmd() *cobra.Command {
	var limit int

	var username string

	cmd := &cobra.Command{
		Use:   "Delete_user_username",
		Short: "Delete user resource.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/user/{username}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{username}", username)

			// build URL and query params
			u := url.URL{Path: pathWithParams}
			q := url.Values{}

			u.RawQuery = q.Encode()
			fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()
			resp, err := http.Get(fullUrl)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var pretty interface{}
			if err := json.Unmarshal(body, &pretty); err != nil {
				return err
			}
			prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(prettyJSON))
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of items")
	cmd.Flags().StringVar(&username, "username", "", "path username parameter")

	return cmd
}
