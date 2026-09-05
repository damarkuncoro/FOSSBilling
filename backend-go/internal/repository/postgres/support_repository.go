package postgres

import (
	"context"
	"errors"

	"github.com/damarkuncoro/FOSSBilling/backend-go/internal/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupportRepository struct {
	pool *pgxpool.Pool
}

func NewSupportRepository(pool *pgxpool.Pool) *SupportRepository {
	return &SupportRepository{pool: pool}
}

func (r *SupportRepository) GetTicketByID(ctx context.Context, id int64) (*domain.Ticket, error) {
	query := `
		SELECT id, client_id, helpdesk_id, subject, status, priority, rel_type, rel_id, created_at, updated_at
		FROM support_tickets
		WHERE id = $1
	`
	var t domain.Ticket
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.ClientID, &t.HelpdeskID, &t.Subject, &t.Status, &t.Priority, &t.RelType, &t.RelID, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appErrors.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *SupportRepository) ListTicketsByClientID(ctx context.Context, clientID int64, limit, offset int) ([]*domain.Ticket, int, error) {
	countQuery := `SELECT COUNT(*) FROM support_tickets WHERE client_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, clientID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, client_id, helpdesk_id, subject, status, priority, rel_type, rel_id, created_at, updated_at
		FROM support_tickets
		WHERE client_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tickets []*domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		if err := rows.Scan(
			&t.ID, &t.ClientID, &t.HelpdeskID, &t.Subject, &t.Status, &t.Priority, &t.RelType, &t.RelID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		tickets = append(tickets, &t)
	}

	return tickets, total, nil
}

func (r *SupportRepository) ListTickets(ctx context.Context, limit, offset int) ([]*domain.Ticket, int, error) {
	countQuery := `SELECT COUNT(*) FROM support_tickets`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, client_id, helpdesk_id, subject, status, priority, rel_type, rel_id, created_at, updated_at
		FROM support_tickets
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tickets []*domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		if err := rows.Scan(
			&t.ID, &t.ClientID, &t.HelpdeskID, &t.Subject, &t.Status, &t.Priority, &t.RelType, &t.RelID, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		tickets = append(tickets, &t)
	}

	return tickets, total, nil
}

func (r *SupportRepository) CreateTicket(ctx context.Context, ticket *domain.Ticket, initialMessage *domain.TicketMessage) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if ticket.Status == "" {
		ticket.Status = domain.TicketStatusOpen
	}
	if ticket.Priority == "" {
		ticket.Priority = domain.PriorityMedium
	}

	ticketQuery := `
		INSERT INTO support_tickets (client_id, helpdesk_id, subject, status, priority, rel_type, rel_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`
	if err := tx.QueryRow(ctx, ticketQuery,
		ticket.ClientID, ticket.HelpdeskID, ticket.Subject, ticket.Status, ticket.Priority, ticket.RelType, ticket.RelID,
	).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt); err != nil {
		return err
	}

	if initialMessage != nil {
		initialMessage.TicketID = ticket.ID
		msgQuery := `
			INSERT INTO support_ticket_messages (ticket_id, admin_id, client_id, content, ip_address, created_at)
			VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
			RETURNING id, created_at
		`
		if err := tx.QueryRow(ctx, msgQuery,
			initialMessage.TicketID, initialMessage.AdminID, initialMessage.ClientID, initialMessage.Content, initialMessage.IPAddress,
		).Scan(&initialMessage.ID, &initialMessage.CreatedAt); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *SupportRepository) AddMessage(ctx context.Context, msg *domain.TicketMessage) error {
	query := `
		INSERT INTO support_ticket_messages (ticket_id, admin_id, client_id, content, ip_address, created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, created_at
	`
	return r.pool.QueryRow(ctx, query,
		msg.TicketID, msg.AdminID, msg.ClientID, msg.Content, msg.IPAddress,
	).Scan(&msg.ID, &msg.CreatedAt)
}

func (r *SupportRepository) UpdateTicketStatus(ctx context.Context, id int64, status domain.TicketStatus) error {
	query := `UPDATE support_tickets SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	tag, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return appErrors.ErrNotFound
	}
	return nil
}

func (r *SupportRepository) GetMessages(ctx context.Context, ticketID int64) ([]*domain.TicketMessage, error) {
	query := `
		SELECT id, ticket_id, admin_id, client_id, content, ip_address, created_at
		FROM support_ticket_messages
		WHERE ticket_id = $1
		ORDER BY id ASC
	`
	rows, err := r.pool.Query(ctx, query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.TicketMessage
	for rows.Next() {
		var m domain.TicketMessage
		if err := rows.Scan(
			&m.ID, &m.TicketID, &m.AdminID, &m.ClientID, &m.Content, &m.IPAddress, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}
	return messages, nil
}
