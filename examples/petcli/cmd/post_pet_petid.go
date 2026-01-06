package cmd

import (
	"net/url"
	"strings"
	"strconv"
	"github.com/spf13/cobra"
	"io/ioutil"
	"encoding/json"
	"petcli/config"
	"fmt"
	"net/http"

)
func NewPost_pet_petIdCmd() *cobra.Command {
	var limit int

	var petId int
	var name string
	var status string

	cmd := &cobra.Command{
		Use:   "Post_pet_petId",
		Short: "Updates a pet in the store with form data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/pet/{petId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{petId}", strconv.Itoa(petId))

			// build URL and query params
			u := url.URL{Path: pathWithParams}
			q := url.Values{}
	if name != "" { q.Set("name", name) }
	if status != "" { q.Set("status", status) }

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
	cmd.Flags().IntVar(&petId, "petId", 0, "path petId parameter")
	cmd.Flags().StringVar(&name, "name", "", "query name parameter")
	cmd.Flags().StringVar(&status, "status", "", "query status parameter")

	return cmd
}
