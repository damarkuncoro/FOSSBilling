package massmail

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/mailer"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

type Service interface {
	Create(ctx context.Context, adminID int64, subject, content string) (*domain.MassMailCampaign, error)
	GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error)
	List(ctx context.Context, limit, offset int) ([]*domain.MassMailCampaign, int, error)
	Send(ctx context.Context, campaignID int64) (*domain.MassMailCampaign, error)
}

type MassMailService struct {
	repo domain.MassMailRepository; clientRepo domain.ClientRepository; mailer mailer.Mailer; from, app string
}

func NewMassMailService(r domain.MassMailRepository, cr domain.ClientRepository, m mailer.Mailer, from, app string) *MassMailService {
	return &MassMailService{r, cr, m, from, app}
}

func (s *MassMailService) Create(ctx context.Context, aID int64, sub, cont string) (*domain.MassMailCampaign, error) {
	if sub == "" || cont == "" { return nil, appErrors.ErrInvalidInput }
	c := &domain.MassMailCampaign{AdminID: aID, Subject: security.SanitizeAlphaNumeric(sub), Content: security.SanitizeHTML(cont), Status: domain.CampaignStatusDraft}
	return c, s.repo.Create(ctx, c)
}

func (s *MassMailService) Send(ctx context.Context, id int64) (*domain.MassMailCampaign, error) {
	cp, err := s.repo.GetByID(ctx, id)
	if err != nil { return nil, err }
	if cp.Status == domain.CampaignStatusCompleted { return nil, errors.New("already completed") }
	cls, _, _ := s.clientRepo.List(ctx, 5000, 0)
	cp.Status = domain.CampaignStatusSending; _ = s.repo.Update(ctx, cp)

	sent := 0
	for _, c := range cls {
		if c.Email != "" && s.mailer.Send(ctx, mailer.Message{From: fmt.Sprintf("%s <%s>", s.app, s.from), To: []string{c.Email}, Subject: cp.Subject, HTMLBody: cp.Content}) == nil {
			sent++
		}
	}
	cp.SentCount, cp.Status, cp.SentAt = sent, domain.CampaignStatusCompleted, pointer(time.Now().UTC())
	return cp, s.repo.Update(ctx, cp)
}

func (s *MassMailService) GetByID(ctx context.Context, id int64) (*domain.MassMailCampaign, error) { return s.repo.GetByID(ctx, id) }

func (s *MassMailService) List(ctx context.Context, l, o int) ([]*domain.MassMailCampaign, int, error) {
	if l <= 0 { l = 20 }; return s.repo.List(ctx, l, o)
}

func pointer[T any](v T) *T { return &v }
