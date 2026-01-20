package cmd

import (
	"io"
	"strings"
	"petcli_models/utils"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"petcli_models/config"
	"strconv"
)

func NewGetStoreOrderOrderIdCmd() *cobra.Command {
	var orderIdPath int

    cmd := &cobra.Command{
        Use:   "get-store-order-orderid",
        Short: "Find purchase order by ID.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/store/order/{orderId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{orderId}", strconv.Itoa(orderIdPath))

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
	cmd.Flags().IntVar(&orderIdPath, "orderId", 0, "path orderId parameter")
	cmd.MarkFlagRequired("orderId")

    return cmd
}
