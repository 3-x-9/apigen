package cmd

import (
	"net/http"
	"io"
	"encoding/json"
	"net/url"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"petcli/config"

)
func NewGet_Pet_FindByTagsCmd() *cobra.Command {
	var tags string

    cmd := &cobra.Command{
        Use:   "Get_Pet_FindByTags",
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

			var bodyReader io.Reader = nil

            req, err := http.NewRequest("GET", fullUrl, bodyReader)
            if err != nil {
                return err
            }

			
		if cfg.ApiKey != "" {
			req.Header.Set("api_key", cfg.ApiKey)}

			if Debug {
				fmt.Println("---DEBUG INFO---")
				fmt.Printf("%-15s: %s\n", "Request Method", req.Method)
				fmt.Printf("%-15s: %s\n", "URL", req.URL.String())
				fmt.Printf("%-15s: %v\n", "Headers", req.Header)
				fmt.Println("----------------------")
			}

			

            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            body, err := io.ReadAll(resp.Body)
            if err != nil {
                return err
            }
            var pretty interface{}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				fmt.Println("Request failed:")
				fmt.Printf("%-15s: %s\n", "Error", resp.Status)
				fmt.Printf("%-15s: %s\n", "URL", resp.Request.URL.String())
				fmt.Printf("%-15s: %s\n", "METHOD", resp.Request.Method)
				fmt.Println("----------------------")
}
			if strings.Contains(resp.Header.Get("Content-Type"), "json") {
            if err := json.Unmarshal(body, &pretty); err != nil {
                return err
            }
            prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
            if err != nil {
                return err
            }
            fmt.Println("Response body:\n" + string(prettyJSON))
			} else {
			 	fmt.Println("Response body:\n" + string(body))
			}
            return nil
        },
    }
	cmd.Flags().StringVar(&tags, "tags", "", "query tags parameter")
	cmd.MarkFlagRequired("tags")

    return cmd
}
