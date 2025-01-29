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

type ValidateCollectionCommand struct{}

// Execute runs the validation command
func (v *ValidateCollectionCommand) Execute(params map[string]interface{}) error {
	// Retrieve parameters
	apiEndpoint, apiKey, expectedSubstring, collectionID, err := v.retrieveParams(params)
	if err != nil {
		return err
	}

	// Validate collection information
	err = v.validateCollectionInfo(apiEndpoint, apiKey, expectedSubstring, collectionID)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// retrieveParams processes the input params map and returns the required parameters
func (t *ValidateCollectionCommand) retrieveParams(params map[string]interface{}) (string, string, string, int, error) {

	// Retrieve endpoint url
	apiEndpoint, ok := params["endpoint_url"].(string)
	if !ok {
		return "", "", "", 0, fmt.Errorf("endpoint URL is required")
	}

	// Retrieve API key
	apiKey, ok := params["api_key"].(string)
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

	return apiEndpoint, apiKey, randomString, collectionIDInt, nil
}

// validateCollectionInfo validates the collection name and description
func (cmd *ValidateCollectionCommand) validateCollectionInfo(apiEndpoint, apiKey, expectedSubstring string, collectionID int) error {
	// Build the API URL
	url := fmt.Sprintf("%s/api/collection/%d", apiEndpoint, collectionID)

	// Fetch the collection data
	resp, err := makeHTTPRequest(url, apiKey, getMethod, nil)
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

	// Validate fields
	if err := validateFieldContainsSubstring(collection["name"], expectedSubstring, "name"); err != nil {
		return err
	}

	if err := validateFieldContainsSubstring(collection["description"], expectedSubstring, "description"); err != nil {
		return err
	}

	fmt.Println("Collection name and description are valid")
	return nil
}

// validateFieldContainsSubstring ensures the field contains the expected substring
func validateFieldContainsSubstring(field interface{}, substring, fieldName string) error {
	// Check if the field is a string
	fieldStr, ok := field.(string)
	if !ok {
		return fmt.Errorf("field '%s' is not a valid string", fieldName)
	}

	// Check if the string contains the expected substring
	if !strings.Contains(fieldStr, substring) {
		return fmt.Errorf("field '%s' does not contain the expected substring: '%s'", fieldName, substring)
	}

	return nil
}
