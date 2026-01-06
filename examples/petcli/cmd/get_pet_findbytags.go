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
func NewGet_pet_findByTagsCmd() *cobra.Command {
	var limit int

	var tags string

	cmd := &cobra.Command{
		Use:   "Get_pet_findByTags",
		Short: "Finds Pets by tags.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/pet/findByTags"

			// build URL and query params
			u := url.URL{Path: pathWithParams}
			q := url.Values{}
	if tags != "" { q.Set("tags", tags) }

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
	cmd.Flags().StringVar(&tags, "tags", "", "query tags parameter")

	return cmd
}
