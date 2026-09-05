package support

import (
	"context"
	"errors"
	"strings"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
)

var (
	ErrEmptyTicketSubject = errors.New("ticket subject cannot be empty")
	ErrEmptyTicketMessage = errors.New("ticket message cannot be empty")
	ErrTicketClosed       = errors.New("cannot reply to a closed ticket")
)

type CreateTicketDTO struct {
	ClientID   int64                 `json:"client_id"`
	HelpdeskID int64                 `json:"helpdesk_id"`
	Subject    string                `json:"subject"`
	Message    string                `json:"message"`
	Priority   domain.TicketPriority `json:"priority"`
	RelType    *string               `json:"rel_type,omitempty"`
	RelID      *int64                `json:"rel_id,omitempty"`
	IPAddress  string                `json:"ip_address,omitempty"`
}

type SupportService struct {
	supportRepo domain.SupportRepository
	clientRepo  domain.ClientRepository
}

func NewSupportService(
	supportRepo domain.SupportRepository,
	clientRepo domain.ClientRepository,
) *SupportService {
	return &SupportService{
		supportRepo: supportRepo,
		clientRepo:  clientRepo,
	}
}

// OpenTicket creates a new support ticket with initial message
func (s *SupportService) OpenTicket(ctx context.Context, dto CreateTicketDTO) (*domain.Ticket, error) {
	if strings.TrimSpace(dto.Subject) == "" {
		return nil, ErrEmptyTicketSubject
	}
	if strings.TrimSpace(dto.Message) == "" {
		return nil, ErrEmptyTicketMessage
	}

	if _, err := s.clientRepo.GetByID(ctx, dto.ClientID); err != nil {
		return nil, err
	}

	if dto.Priority == "" {
		dto.Priority = domain.PriorityMedium
	}

	ticket := &domain.Ticket{
		ClientID:   dto.ClientID,
		HelpdeskID: dto.HelpdeskID,
		Subject:    dto.Subject,
		Status:     domain.TicketStatusOpen,
		Priority:   dto.Priority,
		RelType:    dto.RelType,
		RelID:      dto.RelID,
	}

	initialMsg := &domain.TicketMessage{
		ClientID:  &dto.ClientID,
		Content:   dto.Message,
		IPAddress: dto.IPAddress,
	}

	if err := s.supportRepo.CreateTicket(ctx, ticket, initialMsg); err != nil {
		return nil, err
	}

	return s.supportRepo.GetTicketByID(ctx, ticket.ID)
}

// ClientReply posts a reply from client and changes status to awaiting_staff
func (s *SupportService) ClientReply(ctx context.Context, ticketID, clientID int64, message, ipAddress string) (*domain.TicketMessage, error) {
	if strings.TrimSpace(message) == "" {
		return nil, ErrEmptyTicketMessage
	}

	ticket, err := s.supportRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.ClientID != clientID {
		return nil, appErrors.ErrForbidden
	}

	if ticket.Status == domain.TicketStatusClosed {
		return nil, ErrTicketClosed
	}

	msg := &domain.TicketMessage{
		TicketID:  ticketID,
		ClientID:  &clientID,
		Content:   message,
		IPAddress: ipAddress,
	}

	if err := s.supportRepo.AddMessage(ctx, msg); err != nil {
		return nil, err
	}

	_ = s.supportRepo.UpdateTicketStatus(ctx, ticketID, domain.TicketStatusAwaitingStaff)

	return msg, nil
}

// StaffReply posts a reply from staff and changes status to awaiting_client
func (s *SupportService) StaffReply(ctx context.Context, ticketID, adminID int64, message string) (*domain.TicketMessage, error) {
	if strings.TrimSpace(message) == "" {
		return nil, ErrEmptyTicketMessage
	}

	ticket, err := s.supportRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.Status == domain.TicketStatusClosed {
		return nil, ErrTicketClosed
	}

	msg := &domain.TicketMessage{
		TicketID: ticketID,
		AdminID:  &adminID,
		Content:  message,
	}

	if err := s.supportRepo.AddMessage(ctx, msg); err != nil {
		return nil, err
	}

	_ = s.supportRepo.UpdateTicketStatus(ctx, ticketID, domain.TicketStatusAwaitingClient)

	return msg, nil
}

// CloseTicket sets status to closed
func (s *SupportService) CloseTicket(ctx context.Context, ticketID int64) error {
	return s.supportRepo.UpdateTicketStatus(ctx, ticketID, domain.TicketStatusClosed)
}
