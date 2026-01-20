package cmd

import (
	"petcli/config"
	"petcli/utils"
	"github.com/spf13/cobra"
	"net/http"
	"io"
	"net/url"
	"strings"
)

func NewPutUserUsernameCmd() *cobra.Command {
	var lastName string
	var password string
	var phone string
	var userStatus int
	var username string
	var email string
	var firstName string
	var id int
	var body string
	var contentType string
	var usernamePath string

    cmd := &cobra.Command{
        Use:   "put-user-username",
        Short: "Update user resource.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli", Env)
            pathWithParams := "/user/{username}"
	pathWithParams = strings.ReplaceAll(pathWithParams, "{username}", usernamePath)

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
		
	if lastName != "" {
		bodyObj["lastName"] = lastName
		}
		
	if password != "" {
		bodyObj["password"] = password
		}
		
	if phone != "" {
		bodyObj["phone"] = phone
		}
		
	if userStatus != 0 {
		bodyObj["userStatus"] = userStatus
		}
		
	if username != "" {
		bodyObj["username"] = username
		}
		
	if email != "" {
		bodyObj["email"] = email
		}
		
	if firstName != "" {
		bodyObj["firstName"] = firstName
		}
		
	if id != 0 {
		bodyObj["id"] = id
		}
		if err := utils.CheckBody(body, bodyObj, &bodyReader); err != nil {
    return err
}
			if body == "" {
				
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
	cmd.Flags().StringVar(&lastName, "body-lastName", "", "lastName parameter")
	cmd.Flags().StringVar(&password, "body-password", "", "password parameter")
	cmd.Flags().StringVar(&phone, "body-phone", "", "phone parameter")
	cmd.Flags().IntVar(&userStatus, "body-userStatus", 0, "User Status")
	cmd.Flags().StringVar(&username, "body-username", "", "username parameter")
	cmd.Flags().StringVar(&email, "body-email", "", "email parameter")
	cmd.Flags().StringVar(&firstName, "body-firstName", "", "firstName parameter")
	cmd.Flags().IntVar(&id, "body-id", 0, "id parameter")
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().StringVar(&contentType, "content-type", "", "Content-Type header for the request body")
	cmd.Flags().StringVar(&usernamePath, "username", "", "path username parameter")
	cmd.MarkFlagRequired("username")

    return cmd
}
