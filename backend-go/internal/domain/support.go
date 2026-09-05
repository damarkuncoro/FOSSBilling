package domain

import (
	"context"
	"time"
)

type TicketStatus string

const (
	TicketStatusOpen             TicketStatus = "open"
	TicketStatusAwaitingClient   TicketStatus = "awaiting_client"
	TicketStatusAwaitingStaff    TicketStatus = "awaiting_staff"
	TicketStatusClosed           TicketStatus = "closed"
)

type TicketPriority string

const (
	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
	PriorityUrgent TicketPriority = "urgent"
)

type Ticket struct {
	ID           int64          `json:"id"`
	ClientID     int64          `json:"client_id"`
	HelpdeskID   int64          `json:"helpdesk_id"`
	Subject      string         `json:"subject"`
	Status       TicketStatus   `json:"status"`
	Priority     TicketPriority `json:"priority"`
	RelType      *string        `json:"rel_type,omitempty"` // e.g. "order"
	RelID        *int64         `json:"rel_id,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type TicketMessage struct {
	ID        int64     `json:"id"`
	TicketID  int64     `json:"ticket_id"`
	AdminID   *int64    `json:"admin_id,omitempty"`
	ClientID  *int64    `json:"client_id,omitempty"`
	Content   string    `json:"content"`
	IPAddress string    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type SupportRepository interface {
	GetTicketByID(ctx context.Context, id int64) (*Ticket, error)
	ListTicketsByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*Ticket, int, error)
	ListTickets(ctx context.Context, limit, offset int) ([]*Ticket, int, error)
	CreateTicket(ctx context.Context, ticket *Ticket, initialMessage *TicketMessage) error
	AddMessage(ctx context.Context, message *TicketMessage) error
	UpdateTicketStatus(ctx context.Context, id int64, status TicketStatus) error
	GetMessages(ctx context.Context, ticketID int64) ([]*TicketMessage, error)
}

