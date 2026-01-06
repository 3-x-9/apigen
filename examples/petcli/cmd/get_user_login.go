package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"encoding/json"
	"net/url"
	"io"
	"io/ioutil"
	"strings"
	"petcli/config"

)
func NewGet_user_loginCmd() *cobra.Command {
	var limit int

	var username string
	var password string

    cmd := &cobra.Command{
        Use:   "Get_user_login",
        Short: "Logs user into the system.",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := config.Load("petcli")
            pathWithParams := "/user/login"

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
	if username != "" { q.Set("username", username) }
	if password != "" { q.Set("password", password) }

            u.RawQuery = q.Encode()
            fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()

			var bodyReader io.Reader = nil

            req, err := http.NewRequest("GET", fullUrl, bodyReader)
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
	cmd.Flags().StringVar(&username, "username", "", "query username parameter")
	cmd.Flags().StringVar(&password, "password", "", "query password parameter")

    return cmd
}
