package provisioning

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
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
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type DirectAdminAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
	Package  string `json:"package"`
	Email    string `json:"email"`
	IP       string `json:"ip"`
	Status   string `json:"status"`
}

func (p *DirectAdminProvisioner) CreateAccount(ctx context.Context, acc DirectAdminAccount) (*DirectAdminAccount, error) {
	if acc.Username == "" {
		acc.Username = "da" + strings.ToLower(fmt.Sprintf("%d", time.Now().UnixNano()%1000000))
	}
	if acc.Password == "" {
		acc.Password = fmt.Sprintf("DA!%s#%d", acc.Username, time.Now().Year())
	}
	if acc.IP == "" {
		acc.IP = p.Host
	}

	acc.Status = "active"
	return &acc, nil
}

func (p *DirectAdminProvisioner) SuspendAccount(ctx context.Context, username, reason string) error {
	// CMD_API_SELECT_USERS suspend command
	return nil
}

func (p *DirectAdminProvisioner) UnsuspendAccount(ctx context.Context, username string) error {
	// CMD_API_SELECT_USERS unsuspend command
	return nil
}
