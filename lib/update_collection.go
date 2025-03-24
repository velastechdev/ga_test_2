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
	"strings"
)

type UpdateCollection struct {
	ApiEndpoint  string
	ApiKey       string
	CollectionID string
	SuffixString string
}

// UpdateCollectionInfo updates the collection name and description
func (u *UpdateCollection) UpdateCollectionInfo() (string, string, error) {
	fmt.Println("Updating collection name and description...")

	// Build the API URL
	url := fmt.Sprintf("%s/api/collection/%s", u.ApiEndpoint, u.CollectionID)

	// Fetch existing collection details
	resp, err := MakeHTTPRequest(url, u.ApiKey, http.MethodGet, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch collection details: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var collection map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Update name field
	name, err := appendSuffixToField(collection["name"], u.SuffixString)
	if err != nil {
		return "", "", err
	}

	// Update description field
	description, err := appendSuffixToField(collection["description"], u.SuffixString)
	if err != nil {
		return "", "", err
	}

	// Prepare updated collection
	updatedCollection := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	// Send the updated collection back to the API
	_, err = MakeHTTPRequest(url, u.ApiKey, http.MethodPut, updatedCollection)
	if err != nil {
		return "", "", fmt.Errorf("failed to update collection: %w", err)
	}

	return name, description, nil
}

// appendSuffixToField appends a suffix to the given field
func appendSuffixToField(field interface{}, suffix string) (string, error) {

	// Ensure the field is a string
	if fieldStr, ok := field.(string); ok {
		if strings.Contains(fieldStr, " | ") {
			// If " | " exists, update the second part
			parts := strings.Split(fieldStr, " | ")
			return fmt.Sprintf("%s | %s", parts[0], suffix), nil
		}
		// If " | " doesn't exist, append the suffix
		return fmt.Sprintf("%s | %s", fieldStr, suffix), nil
	}

	return "", fmt.Errorf("field is not a string: %v", field)
}
