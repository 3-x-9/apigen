package cmd

import (
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"github.com/spf13/cobra"
	"net/http"
	"io"
	"net/url"
)

func NewDeleteUserUsernameCmd() *cobra.Command {
	var usernamePath string

    cmd := &cobra.Command{
        Use:   "delete-user-username",
        Short: "Delete user resource.",
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


            req, err := http.NewRequest("DELETE", fullUrl, bodyReader)
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
