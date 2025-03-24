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
	"strconv"
)

type GetDashboardEntityID struct{}

// Execute runs the export command
func (e *GetDashboardEntityID) Execute(params map[string]interface{}) error {
	// Retrieve parameters
	endpointURL, metabaseApiKey, dashboardID, err := e.retrieveParams(params)
	if err != nil {
		return err
	}

	// Build request URL
	url := fmt.Sprintf("%s/api/dashboard/%d", endpointURL, dashboardID)

	// Make HTTP request
	resp, err := lib.MakeHTTPRequest(url, metabaseApiKey, http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data from Metabase, status code: %d", resp.StatusCode)
	}

	var result Dashboard
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Print the entity ID
	fmt.Printf("ENTITY_ID=%s\n", result.EntityID)

	return nil
}

// retrieveParams processes the input params map and returns the required parameters
func (e *GetDashboardEntityID) retrieveParams(params map[string]interface{}) (string, string, int, error) {

	// Retrieve environment
	endpointURL, ok := params["endpoint_url"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("endpoint URL is required")
	}

	// Retrieve API key
	metabaseApiKey, ok := params["api_key"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("api_key is required")
	}

	// Retrieve collection and determine if exporting all collections
	dashboardID, ok := params["dashboard_id"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("dashboard_id is required")
	}

	dashboardIDInt, err := strconv.Atoi(dashboardID)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to convert dashboard ID to integer: %w", err)
	}

	return endpointURL, metabaseApiKey, dashboardIDInt, nil
}
