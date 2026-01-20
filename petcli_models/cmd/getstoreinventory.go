package cmd

import (
	"net/http"
	"io"
	"net/url"
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"github.com/spf13/cobra"
)

func NewGetStoreInventoryCmd() *cobra.Command {

    cmd := &cobra.Command{
        Use:   "get-store-inventory",
        Short: "Returns pet inventories by status.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/store/inventory"

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




		if cfg.Api_keyAuth != "" {
			req.Header.Set("api_key", cfg.Api_keyAuth)}

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
