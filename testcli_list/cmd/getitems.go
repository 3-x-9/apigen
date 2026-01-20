package cmd

import (
	"strings"
	"testcli_list/config"
	"testcli_list/models"
	"github.com/spf13/cobra"
	"io"
	"net/url"
	"testcli_list/utils"
	"net/http"
)

func NewGetItemsCmd() *cobra.Command {

    cmd := &cobra.Command{
        Use:   "get-items",
        Short: "Get items",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("testcli_list", Env)
            pathWithParams := "/items"

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
            
            
            var respObj []models.Item
            if err := utils.ResponsePrint(resp, &respObj, Output); err != nil {
                return err
            }
            
            return nil
        },
    }

    return cmd
}
