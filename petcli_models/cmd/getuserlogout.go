package cmd

import (
	"petcli_models/config"
	"petcli_models/utils"
	"github.com/spf13/cobra"
	"net/http"
	"io"
	"net/url"
	"strings"
)

func NewGetUserLogoutCmd() *cobra.Command {

    cmd := &cobra.Command{
        Use:   "get-user-logout",
        Short: "Logs out current logged in user session.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/user/logout"

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

    return cmd
}
