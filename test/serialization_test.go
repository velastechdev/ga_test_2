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

package test

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github/hulilabs/hp-metrics-metabase-dashboards/commands"
	"github/hulilabs/hp-metrics-metabase-dashboards/lib"
	"os"
	"testing"
)

// TestSerializationIntegration runs a test to validate the serialization process works that consists
// of the success of validating that a collection name and description were modified from staging to production
func TestSerializationIntegration(t *testing.T) {
	// Retrieve required environment variables
	stgApiEndpoint := mustGetEnv("STAGING_ENDPOINT_URL", t)
	stgApiKey := mustGetEnv("STAGING_API_KEY", t)
	prodApiEndpoint := mustGetEnv("PRODUCTION_ENDPOINT_URL", t)
	prodApiKey := mustGetEnv("PRODUCTION_API_KEY", t)
	// The collection ID used in the Metabase URL
	collectionID := mustGetEnv("TEST_COLLECTION_ID", t)

	// Generate a random string for uniqueness in the test
	randomString, err := generateRandomString()
	if err != nil {
		t.Fatalf("Failed to generate random string: %v", err)
	}

	// Update the collection name and description in the staging environment
	updateCollection := &lib.UpdateCollection{
		ApiEndpoint:  stgApiEndpoint,
		ApiKey:       stgApiKey,
		CollectionID: collectionID,
		SuffixString: randomString,
	}

	updatedName, updatedDescription, err := updateCollection.UpdateCollectionInfo()
	if err != nil {
		t.Fatalf("Failed to update collection ID: %s, error: %v", collectionID, err)
	}

	// Prepare parameters for export and import commands
	paramsExport := createCommandParams(stgApiEndpoint, stgApiKey, collectionID, randomString)
	paramsImport := createCommandParams(prodApiEndpoint, prodApiKey, collectionID, randomString)

	// Create the command executor and register commands
	executor := &commands.CommandExecutor{}
	executor.RegisterCommand("export", &commands.ExportCommand{})
	executor.RegisterCommand("import", &commands.ImportCommand{})

	// Execute the export command
	if err := executor.ExecuteCommand("export", paramsExport); err != nil {
		t.Fatalf("Export command from staging failed: %v", err)
	}

	// Execute the import command
	if err := executor.ExecuteCommand("import", paramsImport); err != nil {
		t.Fatalf("Import command from production failed: %v", err)
	}

	// Validate the collection was updated successfully in the production environment
	validateCollection := &lib.ValidateCollection{
		ApiEndpoint:  prodApiEndpoint,
		ApiKey:       prodApiKey,
		SuffixString: randomString,
		CollectionID: collectionID,
	}

	// Validate that collection name and description were modified from staging to production
	if err := validateCollection.ValidateCollectionInfo(updatedName, updatedDescription); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	t.Log("Integration test completed successfully")
}

// generateRandomString generates a random string of 6 characters.
func generateRandomString() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:6], nil
}

// mustGetEnv retrieves an environment variable or fails the test if it is missing or empty.
func mustGetEnv(key string, t *testing.T) string {
	value := os.Getenv(key)
	if value == "" {
		t.Fatalf("Missing or empty environment variable: %s", key)
	}
	return value
}

// createCommandParams creates a map of parameters for the export/import commands.
func createCommandParams(endpoint, apiKey, collectionID, randomString string) map[string]interface{} {
	return map[string]interface{}{
		"endpoint_url":  endpoint,
		"api_key":       apiKey,
		"random_string": randomString,
		"collection":    collectionID,
	}
}
