package massmail_test

import (
	"context"
	"testing"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/repository/memory"
	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/usecase/massmail"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMassMailService_CRUD_And_Send(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMockMassMailRepository()
	clientRepo := memory.NewMockClientRepository()
	mockMailer := mailer.NewMockMailer()

	// Seed clients
	_ = clientRepo.Create(ctx, &domain.Client{Email: "user1@example.com", FirstName: "Alice"})
	_ = clientRepo.Create(ctx, &domain.Client{Email: "user2@example.com", FirstName: "Bob"})

	service := massmail.NewMassMailService(repo, clientRepo, mockMailer, "admin@fossbilling.org", "FOSSBilling")

	t.Run("Create campaign validation failure", func(t *testing.T) {
		_, err := service.Create(ctx, 1, "", "content")
		assert.Error(t, err)

		_, err = service.Create(ctx, 1, "subject", "")
		assert.Error(t, err)
	})

	t.Run("Create campaign success", func(t *testing.T) {
		campaign, err := service.Create(ctx, 1, "Spring Promo", "<p>Big discounts!</p>")
		require.NoError(t, err)
		assert.NotZero(t, campaign.ID)
		assert.Equal(t, domain.CampaignStatusDraft, campaign.Status)
		assert.Equal(t, "Spring Promo", campaign.Subject)
	})

	t.Run("Get and List campaigns", func(t *testing.T) {
		campaigns, total, err := service.List(ctx, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, campaigns, 1)

		c, err := service.GetByID(ctx, campaigns[0].ID)
		require.NoError(t, err)
		assert.Equal(t, "Spring Promo", c.Subject)
	})

	t.Run("Send campaign success", func(t *testing.T) {
		campaigns, _, err := service.List(ctx, 10, 0)
		require.NoError(t, err)
		require.NotEmpty(t, campaigns)

		cID := campaigns[0].ID
		sentCampaign, err := service.Send(ctx, cID)
		require.NoError(t, err)
		assert.Equal(t, domain.CampaignStatusCompleted, sentCampaign.Status)
		assert.Equal(t, 2, sentCampaign.SentCount)
		assert.NotNil(t, sentCampaign.SentAt)
		assert.Len(t, mockMailer.SentMessages, 2)

		// Trying to send again returns error
		_, err = service.Send(ctx, cID)
		assert.Error(t, err)
	})
}
