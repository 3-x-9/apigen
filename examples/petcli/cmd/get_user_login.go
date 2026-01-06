package cmd

import (
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
	"petcli/config"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"

)
func NewGet_user_loginCmd() *cobra.Command {
	var limit int

	var username string
	var password string

	cmd := &cobra.Command{
		Use:   "Get_user_login",
		Short: "Logs user into the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/user/login"

			// build URL and query params
			u := url.URL{Path: pathWithParams}
			q := url.Values{}
	if username != "" { q.Set("username", username) }
	if password != "" { q.Set("password", password) }

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
	cmd.Flags().StringVar(&username, "username", "", "query username parameter")
	cmd.Flags().StringVar(&password, "password", "", "query password parameter")

	return cmd
}
