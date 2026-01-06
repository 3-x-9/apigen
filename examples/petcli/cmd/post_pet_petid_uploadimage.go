package cmd

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"strings"
	"strconv"
	"fmt"
	"github.com/spf13/cobra"
	"petcli/config"

)
func NewPost_pet_petId_uploadImageCmd() *cobra.Command {
	var limit int

	var petId int
	var additionalMetadata string

	cmd := &cobra.Command{
		Use:   "Post_pet_petId_uploadImage",
		Short: "Uploads an image.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/pet/{petId}/uploadImage"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{petId}", strconv.Itoa(petId))

			// build URL and query params
			u := url.URL{Path: pathWithParams}
			q := url.Values{}
	if additionalMetadata != "" { q.Set("additionalMetadata", additionalMetadata) }

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
	cmd.Flags().StringVar(&additionalMetadata, "additionalMetadata", "", "query additionalMetadata parameter")

	return cmd
}
