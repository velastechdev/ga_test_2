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
package dashboad

import (
	"encoding/json"
	"fmt"
	"github/hulilabs/hp-metrics-metabase-dashboards/lib"
	"net/http"
)

type GetDashboardIDCommand struct{}

// Execute runs the export command
func (e *GetDashboardIDCommand) Execute(params map[string]interface{}) error {
	// Retrieve parameters
	endpointURL, metabaseApiKey, entityID, err := e.retrieveParams(params)
	if err != nil {
		return err
	}

	// Build request URL
	url := fmt.Sprintf("%s/api/dashboard", endpointURL)

	// Make HTTP request
	resp, err := lib.MakeHTTPRequest(url, metabaseApiKey, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data from Metabase, status code: %d", resp.StatusCode)
	}

	var dashboards []Dashboard
	if err := json.NewDecoder(resp.Body).Decode(&dashboards); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	var dashboardID int
	found := false
	for _, dashboard := range dashboards {
		if dashboard.EntityID == entityID {
			dashboardID = dashboard.ID
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("dashboard with entity_id %s not found", entityID)
	}

	// Print the dashboard ID
	fmt.Printf("DASHBOARD_ID:%d\n", dashboardID)

	return nil
}

// retrieveParams processes the input params map and returns the required parameters
func (e *GetDashboardIDCommand) retrieveParams(params map[string]interface{}) (string, string, string, error) {

	// Retrieve environment
	endpointURL, ok := params["endpoint_url"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("endpoint URL is required")
	}

	// Retrieve API key
	metabaseApiKey, ok := params["api_key"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("api_key is required")
	}

	// Retrieve collection and determine if exporting all collections
	entityID, ok := params["entity_id"].(string)
	if !ok {
		return "", "", "", fmt.Errorf("entity_id is required")
	}

	return endpointURL, metabaseApiKey, entityID, nil
}
