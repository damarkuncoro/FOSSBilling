package license

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
)

var (
	ErrNotFound = errors.New("license not found")
	ErrForbidden = errors.New("unauthorized")
)

type LicenseDTO struct {
	ID             int64     `json:"id"`
	ProductName    string    `json:"product_name"`
	LicenseKey     string    `json:"license_key"`
	Status         string    `json:"status"`
	LicensedDomain string    `json:"licensed_domain"`
	LicensedIP     string    `json:"licensed_ip"`
	MaxInstances   int       `json:"max_instances"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type LicenseConfig struct {
	LicenseKey     string `json:"license_key"`
	LicensedDomain string `json:"licensed_domain"`
	LicensedIP     string `json:"licensed_ip"`
	MaxInstances   int    `json:"max_instances"`
}

type LicenseService struct{ repo domain.OrderRepository }

func NewLicenseService(r domain.OrderRepository) *LicenseService { return &LicenseService{r} }

func (s *LicenseService) ListClientLicenses(ctx context.Context, cid int64) ([]LicenseDTO, error) {
	os, _, err := s.repo.ListByClientID(ctx, cid, 100, 0); if err != nil { return nil, err }
	res := make([]LicenseDTO, 0)
	for _, o := range os {
		var c LicenseConfig; _ = json.Unmarshal(o.Config, &c)
		if c.LicenseKey != "" || strings.Contains(strings.ToLower(o.Title), "license") {
			k, m, e := c.LicenseKey, c.MaxInstances, time.Now().AddDate(1, 0, 0)
			if k == "" { k = fmt.Sprintf("FB-ENT-%X", o.ID*48271%65535) }; if m <= 0 { m = 5 }; if o.ExpiresAt != nil { e = *o.ExpiresAt }
			res = append(res, LicenseDTO{ID: o.ID, ProductName: o.Title, LicenseKey: k, Status: string(o.Status), MaxInstances: m, LicensedDomain: c.LicensedDomain, LicensedIP: c.LicensedIP, ExpiresAt: e})
		}
	}
	return res, nil
}

func (s *LicenseService) ResetLicenseLock(ctx context.Context, cid, oid int64) error {
	o, err := s.repo.GetByID(ctx, oid); if err != nil || o.ClientID != cid { return ErrForbidden }
	var c LicenseConfig; _ = json.Unmarshal(o.Config, &c); c.LicensedDomain, c.LicensedIP = "", ""
	o.Config, _ = json.Marshal(c); return s.repo.Update(ctx, o)
}

func (s *LicenseService) ResetLicenseKey(ctx context.Context, cid, oid int64) (string, error) {
	o, err := s.repo.GetByID(ctx, oid); if err != nil || o.ClientID != cid { return "", ErrForbidden }
	var c LicenseConfig; _ = json.Unmarshal(o.Config, &c)
	b := make([]byte, 8); _, _ = rand.Read(b); k := "FB-ENT-" + strings.ToUpper(hex.EncodeToString(b)); c.LicenseKey = k
	o.Config, _ = json.Marshal(c); return k, s.repo.Update(ctx, o)
}
