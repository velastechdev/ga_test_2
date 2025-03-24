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
	"encoding/json"
	"fmt"
	"net/http"
)

type ValidateCollection struct {
	ApiEndpoint  string
	ApiKey       string
	CollectionID string
	SuffixString string
}

// ValidateCollectionInfo validates the collection name and description
func (v *ValidateCollection) ValidateCollectionInfo(updatedName, updatedDescription string) error {
	fmt.Println("Validating collection name and description...")

	// Build the API URL
	url := fmt.Sprintf("%s/api/collection/%s", v.ApiEndpoint, v.CollectionID)

	// Fetch the collection data
	resp, err := MakeHTTPRequest(url, v.ApiKey, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch collection data: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var collection map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if collection["name"] != updatedName {
		return fmt.Errorf("collection name does not match the expected value")
	}

	if collection["description"] != updatedDescription {
		return fmt.Errorf("collection description does not match the expected value")
	}

	fmt.Println("Collection name and description are valid")
	return nil
}
