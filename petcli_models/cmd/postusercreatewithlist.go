package cmd

import (
	"net/url"
	"strings"
	"petcli_models/config"
	"petcli_models/utils"
	"github.com/spf13/cobra"
	"net/http"
	"io"
)

func NewPostUserCreateWithListCmd() *cobra.Command {
	var body string
	var contentType string

    cmd := &cobra.Command{
        Use:   "post-user-createwithlist",
        Short: "Creates list of users with given input array.",
        RunE: func(cmd *cobra.Command, args []string) error {

            cfg := config.Load("petcli_models", Env)
            pathWithParams := "/user/createWithList"

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
		if err := utils.CheckBody(body, bodyObj, &bodyReader); err != nil {
    return err
}
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
	cmd.Flags().StringVarP(&body, "body", "b", "", "Request body (raw JSON, @filename, or '-' for stdin)")
	cmd.Flags().StringVar(&contentType, "content-type", "", "Content-Type header for the request body")

    return cmd
}
