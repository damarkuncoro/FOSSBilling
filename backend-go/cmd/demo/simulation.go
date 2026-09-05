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
			res, _ := cpanelProv.Create(ctx, activated)
			var details map[string]string
			_ = json.Unmarshal(res.AccountDetails, &details)
			fmt.Printf("      📦 cPanel Server   : %s (User: %s)\n", details["server"], details["username"])
		} else if ord.ProductID == 202 {
			daAcc, _ := daProv.CreateAccount(ctx, provisioning.DirectAdminAccount{Domain: "solusinusantara.com", Package: "Business"})
			fmt.Printf("      📦 DirectAdmin     : Host %s (User: %s)\n", daProv.Host, daAcc.Username)
		} else if ord.ProductID == 303 {
			res, _ := licenseProv.Create(ctx, activated)
			var details map[string]string
			_ = json.Unmarshal(res.AccountDetails, &details)
			fmt.Printf("      🔑 Enterprise Key  : %s\n", details["license_key"])
		}
	}

	pleskSub, _ := pleskProv.CreateSubscription(ctx, provisioning.PleskSubscription{DomainName: "plesk-demo.com", PlanName: "Default"})
	fmt.Printf("   🚀 Layanan Plesk   : Domain %s (User: %s)\n", pleskSub.DomainName, pleskSub.Username)
}
