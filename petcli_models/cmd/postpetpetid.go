package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"net/url"
	"strings"
	"petcli_models/utils"
	"net/http"
	"petcli_models/config"
	"strconv"
)

func NewPostPetPetIdCmd() *cobra.Command {
	var body string
	var contentType string
	var petIdPath int
	var name string
	var status string

    cmd := &cobra.Command{
        Use:   "post-pet-petid",
        Short: "Updates a pet in the store with form data.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/pet/{petId}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{petId}", strconv.Itoa(petIdPath))

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
	if name != "" { q.Set("name", name) }
	if status != "" { q.Set("status", status) }

            u.RawQuery = q.Encode()
            fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()


			// prepare request body (allow raw JSON, @filename, or '-' for stdin)
			var bodyReader io.Reader
			var err error
			if body != "" {
				bodyReader, err = utils.GetBodyReader(body)
				if err != nil {
					return err
				}
			} else {
			 
			if body == "" {
				
				}	
			}


            req, err := http.NewRequest("POST", fullUrl, bodyReader)
            if err != nil {
                return err
            }


			if bodyReader != nil {
			    ct := contentType
				if ct == "" {
					ct = ""}
				req.Header.Set("Content-Type", ct)
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
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().StringVar(&contentType, "content-type", "", "Content-Type header for the request body")
	cmd.Flags().IntVar(&petIdPath, "petId", 0, "path petId parameter")
	cmd.MarkFlagRequired("petId")
	cmd.Flags().StringVar(&name, "name", "", "query name parameter")
	cmd.Flags().StringVar(&status, "status", "", "query status parameter")

    return cmd
}
