package centralalerts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Alert struct {
	ID                    string    `json:"id"`
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
	Type                  string    `json:"type"` // "info", "warning", "danger"
	MinFOSSBillingVersion string    `json:"min_fossbilling_version"`
	MaxFOSSBillingVersion string    `json:"max_fossbilling_version"`
	IncludePreviewBranch  bool      `json:"include_preview_branch"`
	CreatedAt             time.Time `json:"created_at"`
}

type CentralAlertsClient struct {
	apiURL     string
	httpClient *http.Client
}

func NewCentralAlertsClient(apiURL string) *CentralAlertsClient {
	if apiURL == "" {
		apiURL = "https://api.fossbilling.net/central-alerts/v1/"
	}
	return &CentralAlertsClient{
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CentralAlertsClient) FetchAlerts(ctx context.Context) ([]Alert, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiURL+"list", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Fallback empty list gracefully when offline or during airgapped deployments
		return []Alert{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Alert{}, fmt.Errorf("central alerts API returned status: %d", resp.StatusCode)
	}

	var res struct {
		Alerts []Alert `json:"alerts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return []Alert{}, nil
	}

	return res.Alerts, nil
}
