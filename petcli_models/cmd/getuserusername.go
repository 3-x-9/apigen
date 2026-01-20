package cmd

import (
	"net/url"
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"github.com/spf13/cobra"
	"net/http"
	"io"
)

func NewGetUserUsernameCmd() *cobra.Command {
	var usernamePath string

    cmd := &cobra.Command{
        Use:   "get-user-username",
        Short: "Get user by user name.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/user/{username}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{username}", usernamePath)

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}

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
            
            return utils.ResponsePrint(resp)
        },
    }
	cmd.Flags().StringVar(&usernamePath, "username", "", "path username parameter")
	cmd.MarkFlagRequired("username")

    return cmd
}
