package memory

import (
	"context"
	"sync"
	"time"

	"github.com/fossbilling/backend-go/internal/domain"
	appErrors "github.com/fossbilling/backend-go/pkg/errors"
)

type MockSupportRepository struct {
	mu           sync.RWMutex
	tickets      map[int64]*domain.Ticket
	messages     map[int64][]*domain.TicketMessage
	nextTicketID int64
	nextMsgID    int64
}

func NewMockSupportRepository() *MockSupportRepository {
	return &MockSupportRepository{
		tickets:      make(map[int64]*domain.Ticket),
		messages:     make(map[int64][]*domain.TicketMessage),
		nextTicketID: 1,
		nextMsgID:    1,
	}
}

func (r *MockSupportRepository) GetTicketByID(ctx context.Context, id int64) (*domain.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tickets[id]
	if !ok {
		return nil, appErrors.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *MockSupportRepository) ListTicketsByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Ticket, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []*domain.Ticket
	for _, t := range r.tickets {
		if t.ClientID == clientID {
			cp := *t
			matched = append(matched, &cp)
		}
	}

	total := len(matched)
	if offset >= total {
		return []*domain.Ticket{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}
	return matched[offset:end], total, nil
}

func (r *MockSupportRepository) CreateTicket(ctx context.Context, ticket *domain.Ticket, initialMessage *domain.TicketMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	ticket.ID = r.nextTicketID
	r.nextTicketID++
	now := time.Now().UTC()
	ticket.CreatedAt = now
	ticket.UpdatedAt = now
	if ticket.Status == "" {
		ticket.Status = domain.TicketStatusOpen
	}
	if ticket.Priority == "" {
		ticket.Priority = domain.PriorityMedium
	}

	cpTicket := *ticket
	r.tickets[ticket.ID] = &cpTicket

	if initialMessage != nil {
		initialMessage.ID = r.nextMsgID
		r.nextMsgID++
		initialMessage.TicketID = ticket.ID
		initialMessage.CreatedAt = now

		cpMsg := *initialMessage
		r.messages[ticket.ID] = append(r.messages[ticket.ID], &cpMsg)
	}

	return nil
}

func (r *MockSupportRepository) AddMessage(ctx context.Context, msg *domain.TicketMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tickets[msg.TicketID]; !ok {
		return appErrors.ErrNotFound
	}

	msg.ID = r.nextMsgID
	r.nextMsgID++
	msg.CreatedAt = time.Now().UTC()

	cp := *msg
	r.messages[msg.TicketID] = append(r.messages[msg.TicketID], &cp)
	return nil
}

func (r *MockSupportRepository) UpdateTicketStatus(ctx context.Context, id int64, status domain.TicketStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tickets[id]
	if !ok {
		return appErrors.ErrNotFound
	}
	t.Status = status
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (r *MockSupportRepository) GetMessages(ctx context.Context, ticketID int64) ([]*domain.TicketMessage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msgs, ok := r.messages[ticketID]
	if !ok {
		return []*domain.TicketMessage{}, nil
	}
	return msgs, nil
}
