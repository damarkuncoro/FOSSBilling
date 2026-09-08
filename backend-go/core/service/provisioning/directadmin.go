package provisioning

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

type DirectAdminProvisioner struct {
	Host       string
	Port       int
	Username   string
	Password   string
	HTTPClient *http.Client
}

func NewDirectAdminProvisioner(host string, port int, username, password string) *DirectAdminProvisioner {
	if port == 0 {
		port = 2222
	}
	return &DirectAdminProvisioner{
		Host:       host,
		Port:       port,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *DirectAdminProvisioner) Type() domain.ProductType {
	return domain.ProductTypeHosting
}

func (p *DirectAdminProvisioner) call(ctx context.Context, command string, params url.Values) (url.Values, error) {
	apiURL := fmt.Sprintf("https://%s:%d/%s", p.Host, p.Port, command)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(p.Username, p.Password)

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return url.ParseQuery(string(body))
}

func (p *DirectAdminProvisioner) Create(ctx context.Context, order *domain.Order) (*domain.ProvisionResult, error) {
	var accountConfig struct {
		Domain string `json:"domain"`
		Plan   string `json:"plan"`
	}
	if len(order.Config) > 0 {
		_ = json.Unmarshal(order.Config, &accountConfig)
	}
	if accountConfig.Domain == "" {
		return nil, errors.New("domain is required for DirectAdmin provisioning")
	}

	username := "da" + strings.ToLower(fmt.Sprintf("%d", order.ID))
	password := fmt.Sprintf("DA!%s#%d", username, time.Now().Year())

	params := url.Values{}
	params.Set("action", "create")
	params.Set("add", "Submit")
	params.Set("username", username)
	params.Set("email", "client@example.com") // Should be from order/client
	params.Set("passwd", password)
	params.Set("passwd2", password)
	params.Set("domain", accountConfig.Domain)
	params.Set("package", accountConfig.Plan)
	params.Set("ip", p.Host)
	params.Set("notify", "no")

	res, err := p.call(ctx, "CMD_API_ACCOUNT_USER", params)
	if err != nil {
		return nil, err
	}

	if res.Get("error") != "0" {
		return nil, fmt.Errorf("directadmin account creation failed: %s", res.Get("text"))
	}

	details := map[string]string{
		"username": username,
		"password": password,
		"domain":   accountConfig.Domain,
		"server":   p.Host,
		"da_url":   fmt.Sprintf("https://%s:%d", p.Host, p.Port),
	}
	detailsJSON, _ := json.Marshal(details)

	return &domain.ProvisionResult{
		Success:        true,
		RemoteID:       username,
		AccountDetails: detailsJSON,
	}, nil
}

func (p *DirectAdminProvisioner) Suspend(ctx context.Context, order *domain.Order, reason string) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("action", "suspend")
	params.Set("user", username)

	res, err := p.call(ctx, "CMD_API_SELECT_USERS", params)
	if err != nil {
		return err
	}

	if res.Get("error") != "0" {
		return fmt.Errorf("directadmin suspension failed: %s", res.Get("text"))
	}
	return nil
}

func (p *DirectAdminProvisioner) Unsuspend(ctx context.Context, order *domain.Order) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("action", "unsuspend")
	params.Set("user", username)

	res, err := p.call(ctx, "CMD_API_SELECT_USERS", params)
	if err != nil {
		return err
	}

	if res.Get("error") != "0" {
		return fmt.Errorf("directadmin unsuspension failed: %s", res.Get("text"))
	}
	return nil
}

func (p *DirectAdminProvisioner) Renew(ctx context.Context, order *domain.Order) error {
	return nil
}

func (p *DirectAdminProvisioner) Terminate(ctx context.Context, order *domain.Order) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("confirmed", "Confirm")
	params.Set("delete", "yes")
	params.Set("select0", username)

	res, err := p.call(ctx, "CMD_API_SELECT_USERS", params)
	if err != nil {
		return err
	}

	if res.Get("error") != "0" {
		return fmt.Errorf("directadmin termination failed: %s", res.Get("text"))
	}
	return nil
}

func (p *DirectAdminProvisioner) Sync(ctx context.Context, order *domain.Order) (*domain.ServiceStatus, error) {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("user", username)

	res, err := p.call(ctx, "CMD_API_SHOW_USER_CONFIG", params)
	if err != nil {
		return nil, err
	}

	return &domain.ServiceStatus{
		IsActive:    res.Get("suspended") == "OFF",
		RemoteState: "active",
	}, nil
}

func (p *DirectAdminProvisioner) ChangePassword(ctx context.Context, order *domain.Order, newPassword string) error {
	var details map[string]string
	_ = json.Unmarshal(order.Config, &details)
	username := details["username"]

	params := url.Values{}
	params.Set("username", username)
	params.Set("passwd", newPassword)
	params.Set("passwd2", newPassword)

	res, err := p.call(ctx, "CMD_API_USER_PASSWD", params)
	if err != nil {
		return err
	}

	if res.Get("error") != "0" {
		return fmt.Errorf("directadmin password change failed: %s", res.Get("text"))
	}
	return nil
}

func (p *DirectAdminProvisioner) TestConnection(ctx context.Context) error {
	_, err := p.call(ctx, "CMD_API_ADMIN_STATS", url.Values{})
	return err
}
