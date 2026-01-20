package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"io"
	"net/url"
	"petcli_struct/config"
	"petcli_struct/models"
	"strconv"
	"strings"
	"petcli_struct/utils"
)

func NewGetStoreOrderOrderIdCmd() *cobra.Command {
	var orderIdPath int

    cmd := &cobra.Command{
        Use:   "get-store-order-orderid",
        Short: "Find purchase order by ID.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_struct", Env)
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
            
            
            var respObj models.Order
            if err := utils.ResponsePrint(resp, &respObj, Output); err != nil {
                return err
            }
            
            return nil
        },
    }
	cmd.Flags().IntVar(&orderIdPath, "orderId", 0, "path orderId parameter")
	cmd.MarkFlagRequired("orderId")

    return cmd
}
