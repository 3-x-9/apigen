package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"strings"
	"bytes"
	"fmt"
	"io"
	"encoding/json"
	"petcli/config"
	"os"

)
func NewPut_User_UsernameCmd() *cobra.Command {
	var body string
	var username string

    cmd := &cobra.Command{
        Use:   "Put_User_Username",
        Short: "Update user resource.",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := config.Load("petcli")
            pathWithParams := "/user/{username}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{username}", username)

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}

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
						data, err = io.ReadAll(os.Stdin)
						if err != nil {
							return err
						}
					} else {
						data, err = os.ReadFile(fname)
						if err != nil {
							return err
						}
					}
					bodyReader = bytes.NewReader(data)
				} else if body == "-" {
					data, err := io.ReadAll(os.Stdin)
					if err != nil {
						return err
					}
					bodyReader = bytes.NewReader(data)
				} else {
					bodyReader = strings.NewReader(body)
				}
			}

            req, err := http.NewRequest("PUT", fullUrl, bodyReader)
            if err != nil {
                return err
            }

			if bodyReader != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			
		if cfg.ApiKey != "" {
			req.Header.Set("api_key", cfg.ApiKey)}
				
			if Debug {
				fmt.Println("---DEBUG INFO---")
				fmt.Printf("%-15s: %s\n", "Request Method", req.Method)
				fmt.Printf("%-15s: %s\n", "URL", req.URL.String())
				fmt.Printf("%-15s: %v\n", "Headers", req.Header)
				if bodyReader != nil {
					data, err := io.ReadAll(bodyReader)
					if err != nil {
						return err
					}

					var parsed interface{}
					if json.Unmarshal(data, &parsed) == nil {
						prettyDebugJSON, err := json.MarshalIndent(parsed, "", "  ")
						if err != nil {
							return err
						}
					fmt.Printf("Request Body:\n%s\n", prettyDebugJSON)
					} else {
						fmt.Printf("Request Body:\n%s\n", string(data))
					}
					bodyReader = bytes.NewReader(data) // reset bodyReader
					
					} else {
						fmt.Printf("Request Body: (empty)\n")
					}
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
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().StringVar(&username, "username", "", "path username parameter")
	cmd.MarkFlagRequired("username")

    return cmd
}
