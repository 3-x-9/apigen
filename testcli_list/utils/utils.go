package utils
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
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

func ResponsePrint(resp *http.Response, target interface{}, format string) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("DEBUG: ResponsePrint Status=%d Body=%s\n", resp.StatusCode, string(body))

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
			
			if format == "table" {
				TablePrint(target)
				return nil
			}
			if format == "csv" {
				CSVPrint(target)
				return nil
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

func TablePrint(target interface{}) {
	fmt.Printf("DEBUG: TablePrint called with target type=%T\n", target)
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		fmt.Println("Table output only supported for lists")
		return
	}

	if v.Len() == 0 {
		fmt.Println("No data found")
		return
	}

	t := v.Index(0).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Print Headers
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, t.Field(i).Name)
	}
	fmt.Println(strings.Join(headers, "\t"))

	// Print Rows
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		var row []string
		for j := 0; j < item.NumField(); j++ {
			row = append(row, fmt.Sprintf("%v", item.Field(j).Interface()))
		}
		fmt.Println(strings.Join(row, "\t"))
	}
}

func CSVPrint(target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		fmt.Println("CSV output only supported for lists")
		return
	}

	if v.Len() == 0 {
		return
	}

	t := v.Index(0).Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Print Headers
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, t.Field(i).Name)
	}
	fmt.Println(strings.Join(headers, ","))

	// Print Rows
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		var row []string
		for j := 0; j < item.NumField(); j++ {
			val := fmt.Sprintf("%v", item.Field(j).Interface())
			if strings.Contains(val, ",") {
				val = "\"" + val + "\""
			}
			row = append(row, val)
		}
		fmt.Println(strings.Join(row, ","))
	}
}
