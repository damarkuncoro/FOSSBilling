package seed

import (
	"context"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

func SeedClients(ctx context.Context, r *memory.MockClientRepository) {
	if r == nil { return }
	ph, _ := auth.HashPassword("Password123!")
	_ = r.Create(ctx, &domain.Client{ID: 1, Email: "client@fossbilling.org", PasswordHash: ph, FirstName: "Budi", LastName: "Santoso", Company: "PT Solusi Cloud Nusantara", Country: "ID", Currency: "USD", Status: domain.ClientStatusActive})
	_ = r.AddBalanceTransaction(ctx, &domain.ClientBalance{ClientID: 1, Type: "credit", Amount: 5000000, Description: "Initial wallet deposit balance"})
}
