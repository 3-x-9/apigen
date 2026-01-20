package cmd

import (
	"github.com/spf13/cobra"
	"net/url"
	"strings"
	"petcli/config"
	"fmt"
	"net/http"
	"io"
	"petcli/utils"
	"petcli/models"
)

func NewPutPetCmd() *cobra.Command {
	var category_id int
	var category_name string
	var id int
	var name string
	var photoUrls []string
	var status string
	var tags []string
	var body string
	var contentType string

    cmd := &cobra.Command{
        Use:   "put-pet",
        Short: "Update an existing pet.",
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

            cfg := config.Load("petcli", Env)
            pathWithParams := "/pet"

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}

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
			 bodyObj := map[string]interface{}{}
		
	if category_id != 0 {if _, ok := bodyObj["category"]; !ok {
			bodyObj["category"] = map[string]interface{}{}
	}
		
		bodyObj["category"].(map[string]interface{})["id"] = category_id
		}
		
	if category_name != "" {if _, ok := bodyObj["category"]; !ok {
			bodyObj["category"] = map[string]interface{}{}
	}
		
		bodyObj["category"].(map[string]interface{})["name"] = category_name
		}
		
	if id != 0 {
		bodyObj["id"] = id
		}
		
	if name != "" {
		bodyObj["name"] = name
		}
		
	if photoUrls != nil {
		bodyObj["photoUrls"] = photoUrls
		}
		
	if status != "" {
		bodyObj["status"] = status
		}
		
	if tags != nil {
		bodyObj["tags"] = tags
		}
		if err := utils.CheckBody(body, bodyObj, &bodyReader); err != nil {
    return err
}
			if body == "" {
				
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

				}	
			}


            req, err := http.NewRequest("PUT", fullUrl, bodyReader)
            if err != nil {
                return err
            }


			if bodyReader != nil {
			    ct := contentType
				if ct == "" {
					ct = "application/json"}
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
            
            
            var respObj models.Pet
            if err := utils.ResponsePrint(resp, &respObj, Output); err != nil {
                return err
            }
            
            return nil
        },
    }
	cmd.Flags().IntVar(&category_id, "body-category-id", 0, "id parameter")
	cmd.Flags().StringVar(&category_name, "body-category-name", "", "name parameter")
	cmd.Flags().IntVar(&id, "body-id", 0, "id parameter")
	cmd.Flags().StringVar(&name, "body-name", "", "name parameter")
	cmd.MarkFlagRequired("body-name")
	cmd.Flags().StringSliceVar(&photoUrls, "body-photoUrls", nil, "photoUrls parameter")
	cmd.MarkFlagRequired("body-photoUrls")
	cmd.Flags().StringVar(&status, "body-status", "", "pet status in the store (one of: [available pending sold])")
	cmd.Flags().StringSliceVar(&tags, "body-tags", nil, "tags parameter")
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().StringVar(&contentType, "content-type", "", "Content-Type header for the request body")

    return cmd
}
