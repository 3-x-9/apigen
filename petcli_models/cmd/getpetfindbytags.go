package cmd

import (
	"io"
	"net/url"
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
)

func NewGetPetFindByTagsCmd() *cobra.Command {
	var tags []string

    cmd := &cobra.Command{
        Use:   "get-pet-findbytags",
        Short: "Finds Pets by tags.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/pet/findByTags"

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
	if tags != nil { q.Set("tags", strings.Join(func() []string { res := []string{}; for _, v := range tags { res = append(res, fmt.Sprintf("%v", v)) }; return res }(), ",")) }

            u.RawQuery = q.Encode()
            fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()


			var bodyReader io.Reader = nil


            req, err := http.NewRequest("GET", fullUrl, bodyReader)
            if err != nil {
                return err
            }




		if cfg.Petstore_authAuth != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.Petstore_authAuth)
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
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "query tags parameter")
	cmd.MarkFlagRequired("tags")

    return cmd
}
