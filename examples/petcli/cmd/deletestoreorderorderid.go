package cmd

import (
	"petcli/utils"
	"github.com/spf13/cobra"
	"net/http"
	"io"
	"strconv"
	"net/url"
	"strings"
	"petcli/config"
)

func NewDeleteStoreOrderOrderIdCmd() *cobra.Command {
	var orderIdPath int

    cmd := &cobra.Command{
        Use:   "delete-store-order-orderid",
        Short: "Delete purchase order by identifier.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli", Env)
            pathWithParams := "/store/order/{orderId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{orderId}", strconv.Itoa(orderIdPath))

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
            
            
            if err := utils.ResponsePrint(resp, nil, Output); err != nil {
                return err
            }
            
            return nil
        },
    }
	cmd.Flags().IntVar(&orderIdPath, "orderId", 0, "path orderId parameter")
	cmd.MarkFlagRequired("orderId")

    return cmd
}
