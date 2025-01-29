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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"testing"
)

// TestStep represents a single integration test step
type TestStep struct {
	Name   string
	Params map[string]interface{}
}

// getEnv retrieves environment variables and ensures they are not empty
func getEnv(key string, t *testing.T) string {
	value := os.Getenv(key)
	if value == "" {
		t.Fatalf("Integration test failed: missing environment variabsle   %s", key)
	}
	return value
}

// createParams returns a map of test parameters based on API endpoint and key
func createParams(endpoint, apiKey, collectionID, randomString string) map[string]interface{} {
	return map[string]interface{}{
		"endpoint_url":  endpoint,
		"api_key":       apiKey,
		"random_string": randomString,
		"collection":    collectionID,
	}
}

// TestSerializationIntegration runs an integration test for the serialization commands
func TestSerializationIntegration(t *testing.T) {
	// Retrieve required environment variables
	stgApiEndpoint := getEnv("STAGING_ENDPOINT_URL", t)
	stgApiKey := getEnv("STAGING_API_KEY", t)
	// prodApiEndpoint := getEnv("PRODUCTION_ENDPOINT_URL", t)
	// prodApiKey := getEnv("PRODUCTION_API_KEY", t)
	collectionID := getEnv("TEST_COLLECTION_ID", t)

	// Generate a random string
	randomString, err := RandomString()
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	// Define test steps using a struct
	testSteps := []TestStep{
		{"update_collection", createParams(stgApiEndpoint, stgApiKey, collectionID, randomString)},
		{"export", createParams(stgApiEndpoint, stgApiKey, collectionID, randomString)},
		// {"import", createParams(prodApiEndpoint, prodApiKey, collectionID, randomString)},
		{"validate_collection", createParams(stgApiEndpoint, stgApiKey, collectionID, randomString)},
	}

	// Create the command executor and register commands
	executor := &CommandExecutor{}
	executor.RegisterCommand("update_collection", &UpdateCollectionCommand{})
	executor.RegisterCommand("export", &ExportCommand{})
	// executor.RegisterCommand("import", &ImportCommand{})
	executor.RegisterCommand("validate_collection", &ValidateCollectionCommand{})

	// Execute test steps
	for _, step := range testSteps {
		fmt.Printf("Executing step: %s\n", step.Name)
		if err := executor.ExecuteCommand(step.Name, step.Params); err != nil {
			t.Fatalf("Integration test failed at step '%s': %v", step.Name, err)
		}
	}
}

// RandomString retrieves a random string
func RandomString() (string, error) {
	// Creates a random string with a length of 32 characters
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:6], nil
}
