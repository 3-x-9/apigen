package cmd

import (
	"petcli/config"
	"fmt"
	"net/http"
	"io"
	"net/url"
	"strconv"
	"github.com/spf13/cobra"
	"encoding/json"
	"strings"

)
func NewDelete_Pet_PetIdCmd() *cobra.Command {
	var api_key string
	var petId int

    cmd := &cobra.Command{
        Use:   "Delete_Pet_PetId",
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
	cmd.Flags().StringVar(&api_key, "api_key", "", "header api_key parameter")
	cmd.Flags().IntVar(&petId, "petId", 0, "path petId parameter")
	cmd.MarkFlagRequired("petId")

    return cmd
}
