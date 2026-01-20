package cmd

import (
	"net/url"
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"io"
)

func NewPostStoreOrderCmd() *cobra.Command {
	var complete bool
	var id int
	var petId int
	var quantity int
	var shipDate string
	var status string
	var body string
	var contentType string

    cmd := &cobra.Command{
        Use:   "post-store-order",
        Short: "Place an order for a pet.",
        RunE: func(cmd *cobra.Command, args []string) error {

	if fmt.Sprintf("%v", status) != "" {
		valid := false
		allowed := []string{"placed", "approved", "delivered"}
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

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/store/order"

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
		
	if complete != false {
		bodyObj["complete"] = complete
		}
		
	if id != 0 {
		bodyObj["id"] = id
		}
		
	if petId != 0 {
		bodyObj["petId"] = petId
		}
		
	if quantity != 0 {
		bodyObj["quantity"] = quantity
		}
		
	if shipDate != "" {
		bodyObj["shipDate"] = shipDate
		}
		
	if status != "" {
		bodyObj["status"] = status
		}
		if err := utils.CheckBody(body, bodyObj, &bodyReader); err != nil {
    return err
}
			if body == "" {
				
	if fmt.Sprintf("%v", status) != "" {
		valid := false
		allowed := []string{"placed", "approved", "delivered"}
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


            req, err := http.NewRequest("POST", fullUrl, bodyReader)
            if err != nil {
                return err
            }


			if bodyReader != nil {
			    ct := contentType
				if ct == "" {
					ct = "application/json"}
				req.Header.Set("Content-Type", ct)
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
	cmd.Flags().BoolVar(&complete, "body-complete", false, "complete parameter")
	cmd.Flags().IntVar(&id, "body-id", 0, "id parameter")
	cmd.Flags().IntVar(&petId, "body-petId", 0, "petId parameter")
	cmd.Flags().IntVar(&quantity, "body-quantity", 0, "quantity parameter")
	cmd.Flags().StringVar(&shipDate, "body-shipDate", "", "shipDate parameter")
	cmd.Flags().StringVar(&status, "body-status", "", "Order Status (one of: [placed approved delivered])")
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().StringVar(&contentType, "content-type", "", "Content-Type header for the request body")

    return cmd
}
