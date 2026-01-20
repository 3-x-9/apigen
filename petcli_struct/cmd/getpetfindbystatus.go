package cmd

import (
	"net/http"
	"net/url"
	"strings"
	"petcli_struct/config"
	"petcli_struct/models"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"petcli_struct/utils"
)

func NewGetPetFindByStatusCmd() *cobra.Command {
	var status string

    cmd := &cobra.Command{
        Use:   "get-pet-findbystatus",
        Short: "Finds Pets by status.",
        RunE: func(cmd *cobra.Command, args []string) error {

	if fmt.Sprintf("%v", status) != "" {
		valid := false
		allowed := []string{"available", "pending", "sold"}
		for _, a := range allowed {
			if fmt.Sprintf("%v", status) == a {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid value for status: %v (allowed: %v)", status, allowed)
		}
	}

            cfg := config.Load("petcli_struct", Env)
            pathWithParams := "/pet/findByStatus"

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}
	if status != "" { q.Set("status", status) }

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
            
            
            var respObj []models.Pet
            if err := utils.ResponsePrint(resp, &respObj, Output); err != nil {
                return err
            }
            
            return nil
        },
    }
	cmd.Flags().StringVar(&status, "status", "", "query status parameter")
	cmd.MarkFlagRequired("status")

    return cmd
}
