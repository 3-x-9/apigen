package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"net/url"
	"strings"
	"petcli/config"
	"io"
	"strconv"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"

)
func NewPost_pet_petId_uploadImageCmd() *cobra.Command {
	var limit int

	var body string
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

			// prepare request body (allow raw JSON, @filename, or '-' for stdin)
			var bodyReader io.Reader
			if body != "" {
				if strings.HasPrefix(body, "@") {
					fname := strings.TrimPrefix(body, "@")
					var data []byte
					var err error
					if fname == "-" {
						data, err = ioutil.ReadAll(os.Stdin)
						if err != nil {
							return err
						}
					} else {
						data, err = ioutil.ReadFile(fname)
						if err != nil {
							return err
						}
					}
					bodyReader = bytes.NewReader(data)
				} else if body == "-" {
					data, err := ioutil.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					bodyReader = bytes.NewReader(data)
				} else {
					bodyReader = strings.NewReader(body)
				}
			}

            req, err := http.NewRequest("POST", fullUrl, bodyReader)
            if err != nil {
                return err
            }

			if bodyReader != nil {
				req.Header.Set("Content-Type", "application/octet-stream")
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
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().IntVar(&petId, "petId", 0, "path petId parameter")
	cmd.Flags().StringVar(&additionalMetadata, "additionalMetadata", "", "query additionalMetadata parameter")

    return cmd
}
