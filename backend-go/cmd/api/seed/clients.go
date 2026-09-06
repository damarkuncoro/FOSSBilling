package seed

import (
	"context"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
)

// SeedClients populates demo clients and wallet balances
func SeedClients(ctx context.Context, repo *memory.MockClientRepository) {
	if repo == nil {
		return
	}

	passHash, _ := auth.HashPassword("Password123!")

	_ = repo.Create(ctx, &domain.Client{
		ID:           1,
		Email:        "client@fossbilling.org",
		PasswordHash: passHash,
		FirstName:    "Budi",
		LastName:     "Santoso",
		Company:      "PT Solusi Cloud Nusantara",
		Country:      "ID",
		Currency:     "USD",
		Status:       domain.ClientStatusActive,
	})

	_ = repo.AddBalanceTransaction(ctx, &domain.ClientBalance{
		ClientID:    1,
		Type:        "credit",
		Amount:      5000000,
		Description: "Initial wallet deposit balance",
	})
}
