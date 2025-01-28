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
package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type UpdateCollectionCommand struct{}

// Execute runs the import command
func (u *UpdateCollectionCommand) Execute(params map[string]interface{}) error {
	// Retrieve parameters
	apiEndpoint, apiKey, suffixString, collectionID, err := u.retrieveParams(params)
	if err != nil {
		return err
	}

	// Update collection information
	err = u.updateCollectionInfo(apiEndpoint, apiKey, suffixString, collectionID)
	if err != nil {
		return fmt.Errorf("failed to update collection information: %w", err)
	}

	return nil
}

// retrieveParams processes the input params map and returns the required parameters
func (u *UpdateCollectionCommand) retrieveParams(params map[string]interface{}) (string, string, string, int, error) {

	// Retrieve endpoint url
	endpointURL, ok := params["endpoint_url"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("endpoint URL is required")
	}

	// Retrieve API key
	metabaseApiKey, ok := params["api_key"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("api_key is required")
	}

	collectionID, ok := params["collection"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("collection ID is required")
	}

	collectionIDInt, err := strconv.Atoi(collectionID)
	if err != nil {
		return "", "", "", 0, fmt.Errorf("invalid collection ID: '%s'. Must be a valid integer", collectionID)
	}

	randomString, ok := params["random_string"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("random string is required")
	}

	return endpointURL, metabaseApiKey, randomString, collectionIDInt, nil
}

// updateCollectionInfo updates the collection name and description
func (u *UpdateCollectionCommand) updateCollectionInfo(apiEndpoint, apiKey, suffixString string, collectionID int) error {
	fmt.Println("Updating collection name and description...")

	// Build the API URL
	url := fmt.Sprintf("%s/api/collection/%d", apiEndpoint, collectionID)

	// Fetch existing collection details
	resp, err := makeHTTPRequest(url, apiKey, getMethod, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch collection details: %w", err)
	}
	defer resp.Body.Close()

	// Decode the response
	var collection map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Update fields
	name, err := appendSuffixToField(collection["name"], suffixString)
	if err != nil {
		return err
	}

	description, err := appendSuffixToField(collection["description"], suffixString)
	if err != nil {
		return err
	}

	// Prepare updated collection
	updatedCollection := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	// Send the updated collection back to the API
	_, err = makeHTTPRequest(url, apiKey, putMethod, updatedCollection)
	if err != nil {
		return fmt.Errorf("failed to update collection: %w", err)
	}

	return nil
}

// appendSuffixToField appends a suffix to the given field
func appendSuffixToField(field interface{}, suffix string) (string, error) {
	fmt.Printf("Updating field: %v\n", field)

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
