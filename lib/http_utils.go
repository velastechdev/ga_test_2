/**                   _
 *  _             _ _| |_
 * | |           | |_   _|
 * | |___  _   _ | | |_|
 * | '_  \| | | || | | |
 * | | | || |_| || | | |
 * |_| |_|\___,_||_| |_|
 *
 * (c) Huli Inc
 */
package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	serializationOuputFolder = "serialization"
)

// MakeHTTPRequest sends an HTTP request to the specified URL
func MakeHTTPRequest(url, metabaseApiKey, method string, payload interface{}) (*http.Response, error) {
	var requestBody *bytes.Buffer
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to encode payload: %w", err)
		}
		requestBody = bytes.NewBuffer(jsonData)
	} else {
		requestBody = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", metabaseApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return resp, nil
}
