package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/service/provisioning"
)

func runProvisioningDemo(ctx context.Context, orderRepo domain.OrderRepository, orders []*domain.Order, cpanelProv *provisioning.CpanelProvisioner, daProv *provisioning.DirectAdminProvisioner, pleskProv *provisioning.PleskProvisioner, licenseProv *provisioning.LicenseProvisioner) {
	fmt.Println("\n[7] ⚡ Eksekusi Multi-Driver Provisioning Otomatis...")
	for _, ord := range orders {
		activated, _ := orderRepo.GetByID(ctx, ord.ID)
		fmt.Printf("   🚀 Layanan Aktif: %s (Status: %s)\n", activated.Title, activated.Status)

		if ord.ProductID == 101 {
			// Mock Config for cPanel
			activated.Config = []byte(`{"domain":"solusinusantara.com","plan":"Advanced"}`)
			res, err := cpanelProv.Create(ctx, activated)
			details := map[string]string{"server": "sg1.nusantara-cloud.com", "username": "solusinu"}
			if err == nil && res != nil && len(res.AccountDetails) > 0 {
				_ = json.Unmarshal(res.AccountDetails, &details)
			}
			fmt.Printf("      📦 cPanel Server   : %s (User: %s)\n", details["server"], details["username"])
		} else if ord.ProductID == 202 {
			// Mock Config for DirectAdmin
			activated.Config = []byte(`{"domain":"solusinusantara.com","plan":"Business"}`)
			res, err := daProv.Create(ctx, activated)
			details := map[string]string{"server": "da.nusantara-cloud.com", "username": "solusinu"}
			if err == nil && res != nil && len(res.AccountDetails) > 0 {
				_ = json.Unmarshal(res.AccountDetails, &details)
			}
			fmt.Printf("      📦 DirectAdmin     : Host %s (User: %s)\n", details["server"], details["username"])
		} else if ord.ProductID == 303 {
			res, err := licenseProv.Create(ctx, activated)
			details := map[string]string{"license_key": "FOSS-ENT-LIVE-SIMULATION-KEY"}
			if err == nil && res != nil && len(res.AccountDetails) > 0 {
				_ = json.Unmarshal(res.AccountDetails, &details)
			}
			fmt.Printf("      🔑 Enterprise Key  : %s\n", details["license_key"])
		}
	}

	// Mock Config for Plesk
	demoOrd := &domain.Order{ID: 999, ClientID: 1, Config: []byte(`{"domain":"plesk-demo.com"}`)}
	res, err := pleskProv.Create(ctx, demoOrd)
	details := map[string]string{"domain": "plesk-demo.com", "username": "pleskuser"}
	if err == nil && res != nil && len(res.AccountDetails) > 0 {
		_ = json.Unmarshal(res.AccountDetails, &details)
	}
	fmt.Printf("   🚀 Layanan Plesk   : Domain %s (User: %s)\n", details["domain"], details["username"])
}
