package provisioning

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type CloudflareDNSProvider struct {
	apiToken string
	client   *http.Client
}

func NewCloudflareDNSProvider(apiToken string) *CloudflareDNSProvider {
	return &CloudflareDNSProvider{
		apiToken: apiToken,
		client:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (p *CloudflareDNSProvider) call(ctx context.Context, method, path string, body interface{}) (json.RawMessage, error) {
	apiURL := "https://api.cloudflare.com/client/v4" + path

	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var cfResp struct {
		Success bool            `json:"success"`
		Errors  []interface{}   `json:"errors"`
		Result  json.RawMessage `json:"result"`
	}

	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return nil, fmt.Errorf("cloudflare parse error: %w", err)
	}

	if !cfResp.Success {
		return nil, fmt.Errorf("cloudflare error: %s", string(respBody))
	}

	return cfResp.Result, nil
}

func (p *CloudflareDNSProvider) getZoneID(ctx context.Context, domainName string) (string, error) {
	res, err := p.call(ctx, "GET", "/zones?name="+domainName, nil)
	if err != nil {
		return "", err
	}

	var zones []struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(res, &zones)
	if len(zones) == 0 {
		return "", fmt.Errorf("zone not found for %s", domainName)
	}
	return zones[0].ID, nil
}

func (p *CloudflareDNSProvider) ListRecords(ctx context.Context, domainName string) ([]domain.DNSRecord, error) {
	zoneID, err := p.getZoneID(ctx, domainName)
	if err != nil {
		return nil, err
	}

	res, err := p.call(ctx, "GET", "/zones/"+zoneID+"/dns_records", nil)
	if err != nil {
		return nil, err
	}

	var cfRecords []struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Content  string `json:"content"`
		TTL      int    `json:"ttl"`
		Priority int    `json:"priority"`
	}
	_ = json.Unmarshal(res, &cfRecords)

	records := make([]domain.DNSRecord, len(cfRecords))
	for i, r := range cfRecords {
		records[i] = domain.DNSRecord{
			ID:       r.ID,
			Type:     r.Type,
			Name:     r.Name,
			Content:  r.Content,
			TTL:      r.TTL,
			Priority: r.Priority,
		}
	}
	return records, nil
}

func (p *CloudflareDNSProvider) AddRecord(ctx context.Context, domainName string, record domain.DNSRecord) error {
	zoneID, err := p.getZoneID(ctx, domainName)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type":    record.Type,
		"name":    record.Name,
		"content": record.Content,
		"ttl":     record.TTL,
	}
	if record.Priority > 0 {
		payload["priority"] = record.Priority
	}

	_, err = p.call(ctx, "POST", "/zones/"+zoneID+"/dns_records", payload)
	return err
}

func (p *CloudflareDNSProvider) UpdateRecord(ctx context.Context, domainName string, record domain.DNSRecord) error {
	zoneID, err := p.getZoneID(ctx, domainName)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"type":    record.Type,
		"name":    record.Name,
		"content": record.Content,
		"ttl":     record.TTL,
	}
	if record.Priority > 0 {
		payload["priority"] = record.Priority
	}

	_, err = p.call(ctx, "PUT", "/zones/"+zoneID+"/dns_records/"+record.ID, payload)
	return err
}

func (p *CloudflareDNSProvider) DeleteRecord(ctx context.Context, domainName string, recordID string) error {
	zoneID, err := p.getZoneID(ctx, domainName)
	if err != nil {
		return err
	}

	_, err = p.call(ctx, "DELETE", "/zones/"+zoneID+"/dns_records/"+recordID, nil)
	return err
}
