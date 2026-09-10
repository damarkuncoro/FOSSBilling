package support

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/events"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/security"
)

var (
	ErrEmptyTicketSubject = errors.New("ticket subject cannot be empty")
	ErrEmptyTicketMessage = errors.New("ticket message cannot be empty")
	ErrTicketClosed       = errors.New("cannot reply to a closed ticket")
)

type CreateTicketDTO struct {
	ClientID int64; HelpdeskID int64; Subject string; Message string; Priority domain.TicketPriority; RelType *string; RelID *int64; IPAddress string
}

type SupportService struct {
	supportRepo domain.SupportRepository; clientRepo domain.ClientRepository; eventBus *events.EventBus
}

func NewSupportService(sr domain.SupportRepository, cr domain.ClientRepository, eb ...*events.EventBus) *SupportService {
	var bus *events.EventBus
	if len(eb) > 0 { bus = eb[0] }
	return &SupportService{sr, cr, bus}
}

func (s *SupportService) OpenTicket(ctx context.Context, dto CreateTicketDTO) (*domain.Ticket, error) {
	if strings.TrimSpace(dto.Subject) == "" { return nil, ErrEmptyTicketSubject }
	if strings.TrimSpace(dto.Message) == "" { return nil, ErrEmptyTicketMessage }
	if _, err := s.clientRepo.GetByID(ctx, dto.ClientID); err != nil { return nil, err }
	if dto.Priority == "" { dto.Priority = domain.PriorityMedium }

	t := &domain.Ticket{ClientID: dto.ClientID, HelpdeskID: dto.HelpdeskID, Subject: security.SanitizeHTML(dto.Subject), Status: domain.TicketStatusOpen, Priority: dto.Priority, RelType: dto.RelType, RelID: dto.RelID}
	m := &domain.TicketMessage{ClientID: &dto.ClientID, Content: security.SanitizeHTML(dto.Message), IPAddress: dto.IPAddress}
	if err := s.supportRepo.CreateTicket(ctx, t, m); err != nil { return nil, err }

	if s.eventBus != nil { s.eventBus.PublishAsync(ctx, events.Event{Type: events.EventTicketOpened, Payload: domain.TicketOpenedPayload{TicketID: t.ID, ClientID: t.ClientID, Subject: t.Subject, Priority: string(t.Priority), CreatedAt: time.Now().UTC()}}) }
	return s.supportRepo.GetTicketByID(ctx, t.ID)
}

func (s *SupportService) reply(ctx context.Context, tID, cID, aID int64, msg, ip, status, author string) (*domain.TicketMessage, error) {
	if strings.TrimSpace(msg) == "" { return nil, ErrEmptyTicketMessage }
	t, err := s.supportRepo.GetTicketByID(ctx, tID)
	if err != nil { return nil, err }
	if cID > 0 && t.ClientID != cID { return nil, appErrors.ErrForbidden }
	if t.Status == domain.TicketStatusClosed { return nil, ErrTicketClosed }

	m := &domain.TicketMessage{TicketID: tID, Content: security.SanitizeHTML(msg), IPAddress: ip}
	if cID > 0 { m.ClientID = &cID } else if aID > 0 { m.AdminID = &aID }

	if err := s.supportRepo.AddMessage(ctx, m); err != nil { return nil, err }
	_ = s.supportRepo.UpdateTicketStatus(ctx, tID, domain.TicketStatus(status))

	if s.eventBus != nil { s.eventBus.PublishAsync(ctx, events.Event{Type: events.EventTicketReplied, Payload: map[string]interface{}{"ticket_id": tID, "author": author}}) }
	return m, nil
}

func (s *SupportService) ClientReply(ctx context.Context, tID, cID int64, msg, ip string) (*domain.TicketMessage, error) {
	return s.reply(ctx, tID, cID, 0, msg, ip, string(domain.TicketStatusAwaitingStaff), "client")
}

func (s *SupportService) StaffReply(ctx context.Context, tID, aID int64, msg string) (*domain.TicketMessage, error) {
	return s.reply(ctx, tID, 0, aID, msg, "", string(domain.TicketStatusAwaitingClient), "staff")
}

func (s *SupportService) CloseTicket(ctx context.Context, tID, cID int64) error {
	t, err := s.supportRepo.GetTicketByID(ctx, tID)
	if err != nil || (cID > 0 && t.ClientID != cID) { return appErrors.ErrForbidden }
	if err := s.supportRepo.UpdateTicketStatus(ctx, tID, domain.TicketStatusClosed); err != nil { return err }
	if s.eventBus != nil { s.eventBus.PublishAsync(ctx, events.Event{Type: events.EventTicketClosed, Payload: map[string]interface{}{"ticket_id": tID, "closed_by": cID}}) }
	return nil
}

type TicketDetails struct { Ticket *domain.Ticket `json:"ticket"`; Messages []*domain.TicketMessage `json:"messages"` }

func (s *SupportService) GetTicket(ctx context.Context, tID, cID int64) (*TicketDetails, error) {
	t, err := s.supportRepo.GetTicketByID(ctx, tID)
	if err != nil || (cID > 0 && t.ClientID != cID) { return nil, appErrors.ErrForbidden }
	msgs, err := s.supportRepo.GetMessages(ctx, tID)
	return &TicketDetails{Ticket: t, Messages: msgs}, err
}

func (s *SupportService) ListClientTickets(ctx context.Context, cID int64, l, o int) ([]*domain.Ticket, int, error) {
	if l <= 0 { l = 20 }; return s.supportRepo.ListTicketsByClientID(ctx, cID, l, o)
}

func (s *SupportService) ListAllTickets(ctx context.Context, l, o int) ([]*domain.Ticket, int, error) {
	if l <= 0 { l = 20 }; return s.supportRepo.ListTickets(ctx, l, o)
}
