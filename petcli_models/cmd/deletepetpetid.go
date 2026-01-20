package cmd

import (
	"net/http"
	"io"
	"net/url"
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"strconv"
	"github.com/spf13/cobra"
)

func NewDeletePetPetIdCmd() *cobra.Command {
	var api_key string
	var petIdPath int

    cmd := &cobra.Command{
        Use:   "delete-pet-petid",
        Short: "Deletes a pet.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/pet/{petId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{petId}", strconv.Itoa(petIdPath))

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
	cmd.Flags().StringVar(&api_key, "api_key", "", "header api_key parameter")
	cmd.Flags().IntVar(&petIdPath, "petId", 0, "path petId parameter")
	cmd.MarkFlagRequired("petId")

    return cmd
}
