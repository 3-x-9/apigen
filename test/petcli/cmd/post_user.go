package cmd

import (
	"encoding/json"
	"net/url"
	"strings"
	"petcli/config"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"io/ioutil"

)
func NewPost_userCmd() *cobra.Command {
	var limit int


    cmd := &cobra.Command{
        Use:   "Post_user",
        Short: "Create user.",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := config.Load("petcli")
            pathWithParams := "/user"

            // build URL and query params
            u := url.URL{Path: pathWithParams}
            q := url.Values{}

            u.RawQuery = q.Encode()
            fullUrl := strings.TrimRight(cfg.BaseURL, "/") + u.String()
            req, err := http.NewRequest("POST", fullUrl, nil)
            if err != nil {
                return err
            }

            resp, err := http.DefaultClient.Do(req)
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

    return cmd
}
