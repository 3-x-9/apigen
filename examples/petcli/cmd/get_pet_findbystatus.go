package cmd

import (
	"net/url"
	"strings"
	"github.com/spf13/cobra"
	"io"
	"petcli/config"
	"fmt"
	"net/http"
	"encoding/json"

)
func NewGet_Pet_FindByStatusCmd() *cobra.Command {
	var status string

    cmd := &cobra.Command{
        Use:   "Get_Pet_FindByStatus",
        Short: "Finds Pets by status.",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := config.Load("petcli")
            pathWithParams := "/pet/findByStatus"

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
	if status != "" { q.Set("status", status) }

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
	cmd.Flags().StringVar(&status, "status", "", "query status parameter")
	cmd.MarkFlagRequired("status")

    return cmd
}
