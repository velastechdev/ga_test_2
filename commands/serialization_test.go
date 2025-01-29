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

func TestSerializationIntegration(t *testing.T) {
	stgApiEndpoint := os.Getenv("STAGING_ENDPOINT_URL")
	stgApiKey := os.Getenv("STAGING_API_KEY")
	collectionID := os.Getenv("TEST_COLLECTION_ID")

	if stgApiEndpoint == "" || stgApiKey == "" || collectionID == "" {
		t.Fatalf("Integration test failed: missing environment variables")
	}

	randomString, err := RandomString()
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	fmt.Printf("New suffix string: %s\n", randomString)

	params := map[string]interface{}{
		"endpoint_url":  stgApiEndpoint,
		"api_key":       stgApiKey,
		"random_string": randomString,
		"collection":    collectionID,
	}

	command := &UpdateCollectionCommand{}
	err = command.Execute(params)
	if err != nil {
		t.Fatalf("Integration test failed: %v", err)
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
