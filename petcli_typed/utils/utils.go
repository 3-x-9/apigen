package utils
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetBodyReader(body string) (io.Reader, error) {
	var bodyReader io.Reader
	if body != "" {
		if strings.HasPrefix(body, "@") {
			fname := strings.TrimPrefix(body, "@")
			var data []byte
			var err error
			if fname == "-" {
				data, err = io.ReadAll(os.Stdin)
				if err != nil {
					return nil, err
				}
			} else {
				data, err = os.ReadFile(fname)
				if err != nil {
					return nil, err
				}
			}
			bodyReader = bytes.NewReader(data)
		} else if body == "-" {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return nil, err
				}
			bodyReader = bytes.NewReader(data)
		} else {
			bodyReader = bytes.NewReader([]byte(body))
		}
	} else {
		return nil, nil
	}
	return bodyReader, nil
}

func ResponsePrint(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Println("Request failed:")
		fmt.Printf("%-15s: %s\n", "Error", resp.Status)
		fmt.Printf("%-15s: %s\n", "URL", resp.Request.URL.String())
		fmt.Printf("%-15s: %s\n", "METHOD", resp.Request.Method)
		fmt.Println("----------------")
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		if target != nil {
			if err := json.Unmarshal(body, target); err != nil {
				return err
			}
			prettyJSON, err := json.MarshalIndent(target, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println("Response body:\n" + string(prettyJSON))
		} else {
			var pretty interface{}
			if err := json.Unmarshal(body, &pretty); err != nil {
				return err
			}
			prettyJSON, err := json.MarshalIndent(pretty, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println("Response body:\n" + string(prettyJSON))
		}
	} else {
		fmt.Println("Response body:\n" + string(body))
	}
	return nil
}

func DebugPrintRequest(req *http.Request, bodyReader *io.Reader) error {
	fmt.Println("---DEBUG INFO---")
	fmt.Printf("%-15s: %s\n", "Request Method", req.Method)
	fmt.Printf("%-15s: %s\n", "URL", req.URL.String())
	fmt.Printf("%-15s: %v\n", "Headers", req.Header)
	if bodyReader != nil && *bodyReader != nil {
		data, err := io.ReadAll(*bodyReader)
		if err != nil {
			fmt.Println("----------------")
			return err
		}

		var parsed interface{}
		if json.Unmarshal(data, &parsed) == nil {
			prettyDebugJSON, err := json.MarshalIndent(parsed, "", "  ")
			if err != nil {
				fmt.Println("----------------")
				return err
			}
			fmt.Printf("%-15s: %s\n", "Request Body", string(prettyDebugJSON))
		} else {
			fmt.Printf("%-15s: %s\n", "Request Body", string(data))
		}
		*bodyReader = bytes.NewReader(data) // reset bodyReader
	} else {
		fmt.Printf("%-15s: %s\n", "Request Body", "(empty)")
	}
	fmt.Println("----------------")
	return nil
}

func CheckBody(body string, bodyObj map[string]interface{}, bodyReader *io.Reader) error {
	if body == "" && len(bodyObj) <= 0 {
		return fmt.Errorf("request body is required (use either --body or body flags!!!)")
	}
	if body == "" && len(bodyObj) > 0 {
		data, err := json.Marshal(bodyObj)
		if err != nil {
			return err
		}
		*bodyReader = bytes.NewReader(data)
	}
	return nil
}
