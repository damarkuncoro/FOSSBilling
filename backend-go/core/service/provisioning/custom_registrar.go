package provisioning

import (
	"context"
	"time"
)

// CustomRegistrarDriver always responds with success.
// Equivalent to Registrar_Adapter_Custom in PHP.
type CustomRegistrarDriver struct {
	rdapDriver *RDAPRegistrarDriver
}

func NewCustomRegistrarDriver(rdap ...*RDAPRegistrarDriver) *CustomRegistrarDriver {
	var r *RDAPRegistrarDriver
	if len(rdap) > 0 {
		r = rdap[0]
	}
	return &CustomRegistrarDriver{rdapDriver: r}
}

func (d *CustomRegistrarDriver) CheckAvailability(ctx context.Context, domainName string) (*DomainAvailability, error) {
	if d.rdapDriver != nil {
		return d.rdapDriver.CheckAvailability(ctx, domainName)
	}

	return &DomainAvailability{
		DomainName:  domainName,
		IsAvailable: true,
		Price:       0,
		Currency:    "USD",
		CheckedAt:   time.Now().UTC(),
		RegistrarRef: "custom",
	}, nil
}

func (d *CustomRegistrarDriver) RegisterDomain(ctx context.Context, req DomainRegistrationRequest) (*DomainRegistrationResult, error) {
	now := time.Now().UTC()
	return &DomainRegistrationResult{
		DomainName:    req.DomainName,
		Status:        "active",
		RegisteredAt:  now,
		ExpiresAt:     now.AddDate(req.Years, 0, 0),
		Nameservers:   req.Nameservers,
		TransactionID: "CUSTOM-REG-SUCCESS",
	}, nil
}

func (d *CustomRegistrarDriver) RenewDomain(ctx context.Context, domainName string, years int) (*DomainRegistrationResult, error) {
	return &DomainRegistrationResult{
		DomainName:    domainName,
		Status:        "active",
		TransactionID: "CUSTOM-REN-SUCCESS",
	}, nil
}
