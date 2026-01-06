package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"strings"
	"petcli/config"

)
func NewGet_store_order_orderIdCmd() *cobra.Command {
	var limit int

	var orderId int

	cmd := &cobra.Command{
		Use:   "Get_store_order_orderId",
		Short: "Find purchase order by ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.Load("petcli")
			pathWithParams := "/store/order/{orderId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{orderId}", strconv.Itoa(orderId))

			// build URL and query params
			u := url.URL{Path: pathWithParams}
			q := url.Values{}

			u.RawQuery = q.Encode()
			fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()
			resp, err := http.Get(fullUrl)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			var pretty interface{}
			if err := json.Unmarshal(body, &pretty); err != nil {
				return err
			}
			prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(prettyJSON))
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of items")
	cmd.Flags().IntVar(&orderId, "orderId", 0, "path orderId parameter")

	return cmd
}
