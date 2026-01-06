package cmd

import (
	"net/http"
	"io/ioutil"
	"strings"
	"io"
	"strconv"
	"fmt"
	"encoding/json"
	"net/url"
	"petcli/config"
	"github.com/spf13/cobra"

)
func NewDelete_pet_petIdCmd() *cobra.Command {
	var limit int

	var api_key string
	var petId int

    cmd := &cobra.Command{
        Use:   "Delete_pet_petId",
        Short: "Deletes a pet.",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := config.Load("petcli")
            pathWithParams := "/pet/{petId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{petId}", strconv.Itoa(petId))

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
	// header param api_key available in api_key variable

            u.RawQuery = q.Encode()
            fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()

			var bodyReader io.Reader = nil

            req, err := http.NewRequest("DELETE", fullUrl, bodyReader)
            if err != nil {
                return err
            }

			
		if cfg.ApiKey != "" {
			req.Header.Set("api_key", cfg.ApiKey)}
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                return err
            }
            var pretty interface{}

			if strings.Contains(resp.Header.Get("Content-Type"), "json") {
            if err := json.Unmarshal(body, &pretty); err != nil {
                return err
            }
            prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
            if err != nil {
                return err
            }
            fmt.Println(string(prettyJSON))
			} else {
			 	fmt.Println(string(body))
			}
            return nil
        },
    }
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of items")
	cmd.Flags().StringVar(&api_key, "api_key", "", "header api_key parameter")
	cmd.Flags().IntVar(&petId, "petId", 0, "path petId parameter")

    return cmd
}
