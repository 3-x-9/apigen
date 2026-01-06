package cmd

import (
	"strings"
	"petcli/config"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"

)
func NewPut_petCmd() *cobra.Command {
	var limit int


	cmd := &cobra.Command{
		Use:   "Put_pet",
		Short: "Update an existing pet.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/pet"

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

	return cmd
}
