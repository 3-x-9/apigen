package cmd

import (
	"petcli_models/config"
	"github.com/spf13/cobra"
	"strings"
	"petcli_models/utils"
	"strconv"
	"net/http"
	"io"
	"net/url"
)

func NewGetPetPetIdCmd() *cobra.Command {
	var petIdPath int

    cmd := &cobra.Command{
        Use:   "get-pet-petid",
        Short: "Find pet by ID.",
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
	cmd.Flags().IntVar(&petIdPath, "petId", 0, "path petId parameter")
	cmd.MarkFlagRequired("petId")

    return cmd
}
