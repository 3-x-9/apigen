package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"io"
	"net/url"
	"strings"
	"petcli/config"
	"petcli/utils"
)

func NewGetUserLoginCmd() *cobra.Command {
	var username string
	var password string

    cmd := &cobra.Command{
        Use:   "get-user-login",
        Short: "Logs user into the system.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli", Env)
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





			if Debug {
				utils.DebugPrintRequest(req, &bodyReader)
			}

            client := &http.Client{
                Timeout: cfg.Timeout,
            }
            resp, err := client.Do(req)
            if err != nil {
                return err
            }
            
            
            if err := utils.ResponsePrint(resp, nil, Output); err != nil {
                return err
            }
            
            return nil
        },
    }
	cmd.Flags().StringVar(&username, "username", "", "query username parameter")
	cmd.Flags().StringVar(&password, "password", "", "query password parameter")

    return cmd
}
