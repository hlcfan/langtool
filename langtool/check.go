package langtool

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	baseURL   = "https://languagetool.org/api"
	checkPath = "/v2/check"
)

var errEmptyInput = errors.New("empty input file")

func Check() error {
	inputFileName := flag.String("f", "", "input file")
	flag.Parse()

	if len(*inputFileName) == 0 {
		return errEmptyInput
	}

	inputFile, errOpen := os.Open(*inputFileName)
	if errOpen != nil {
		return errOpen
	}
	defer inputFile.Close()

	file, errRead := io.ReadAll(inputFile)
	if errRead != nil {
		return fmt.Errorf("failed to read file, error: %w", errRead)
	}

	data := url.Values{}
	data.Set("text", string(file))
	data.Set("language", "en-US")
	data.Set("enabledOnly", "false")
	req, err := http.NewRequest("POST", baseURL+checkPath, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request, error: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, errPost := client.Do(req)
	if errPost != nil {
		return fmt.Errorf("failed to call languagetool.org, error: %w", errPost)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to query languagetool, http code: %d\n", resp.StatusCode)
	}

	var result Result
	dec := json.NewDecoder(resp.Body)
	errDecode := dec.Decode(&result)
	if errDecode != nil {
		return fmt.Errorf("failed to decode response body, error: %w", errDecode)
	}

	for _, match := range result.Matches {
		text := match.Context.Text
		var newLine strings.Builder
		for i, r := range text {
			if i == match.Context.Offset {
				newLine.WriteString("\u001b[4m")
			}

			if i == match.Context.Offset+match.Context.Length {
				newLine.WriteString("\u001b[m")
			}

			fmt.Fprintf(&newLine, "%c", r)
		}

		fmt.Printf("%s\n", newLine.String())

		space := strings.Repeat(" ", match.Context.Offset)
		fmt.Printf("%v👆 %v.\n", space, match.Message)
	}

	return nil
}
