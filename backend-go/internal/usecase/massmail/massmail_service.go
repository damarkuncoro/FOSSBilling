package massmail

import (
	"context"
	"fmt"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
	"github.com/fossbilling/backend-go/pkg/mailer"
)

type Service interface {
	Create(ctx context.Context, adminID int64, subject, content string) (*domain.MassMailCampaign, error)
	GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error)
	List(ctx context.Context, limit, offset int) ([]*domain.MassMailCampaign, int, error)
	Send(ctx context.Context, campaignID int64) (*domain.MassMailCampaign, error)
}

type MassMailService struct {
	repo       domain.MassMailRepository
	clientRepo domain.ClientRepository
	mailer     mailer.Mailer
	fromEmail  string
	appName    string
}

func NewMassMailService(
	repo domain.MassMailRepository,
	clientRepo domain.ClientRepository,
	m mailer.Mailer,
	fromEmail, appName string,
) *MassMailService {
	if fromEmail == "" {
		fromEmail = "no-reply@fossbilling.org"
	}
	if appName == "" {
		appName = "FOSSBilling"
	}
	return &MassMailService{
		repo:       repo,
		clientRepo: clientRepo,
		mailer:     m,
		fromEmail:  fromEmail,
		appName:    appName,
	}
}

func (s *MassMailService) Create(ctx context.Context, adminID int64, subject, content string) (*domain.MassMailCampaign, error) {
	if subject == "" {
		return nil, fmt.Errorf("%w: subject is required", appErrors.ErrInvalidInput)
	}
	if content == "" {
		return nil, fmt.Errorf("%w: content is required", appErrors.ErrInvalidInput)
	}

	campaign := &domain.MassMailCampaign{
		AdminID: adminID,
		Subject: subject,
		Content: content,
		Status:  domain.CampaignStatusDraft,
	}

	if err := s.repo.Create(ctx, campaign); err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *MassMailService) GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MassMailService) List(ctx context.Context, limit, offset int) ([]*domain.MassMailCampaign, int, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *MassMailService) Send(ctx context.Context, campaignID int64) (*domain.MassMailCampaign, error) {
	campaign, err := s.repo.GetByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}

	if campaign.Status == domain.CampaignStatusCompleted {
		return nil, fmt.Errorf("%w: campaign already completed", appErrors.ErrInvalidInput)
	}

	// Fetch all active clients
	clients, _, err := s.clientRepo.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	campaign.Status = domain.CampaignStatusSending
	_ = s.repo.Update(ctx, campaign)

	sentCount := 0
	for _, c := range clients {
		if c.Email == "" {
			continue
		}
		msg := mailer.Message{
			From:     fmt.Sprintf("%s <%s>", s.appName, s.fromEmail),
			To:       []string{c.Email},
			Subject:  campaign.Subject,
			HTMLBody: campaign.Content,
		}
		if err := s.mailer.Send(ctx, msg); err == nil {
			sentCount++
		}
	}

	now := time.Now().UTC()
	campaign.SentCount = sentCount
	campaign.Status = domain.CampaignStatusCompleted
	campaign.SentAt = &now

	if err := s.repo.Update(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}
