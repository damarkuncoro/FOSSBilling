package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	appErrors "github.com/damarkuncoro/FOSSBilling/backend-go/pkg/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ticketCols = `id, client_id, helpdesk_id, subject, status, priority, rel_type, rel_id, created_at, updated_at`

type SupportRepository struct{ pool *pgxpool.Pool }

func NewSupportRepository(p *pgxpool.Pool) *SupportRepository { return &SupportRepository{p} }

func scanTicket(r pgx.Row) (*domain.Ticket, error) {
	var t domain.Ticket; err := r.Scan(&t.ID, &t.ClientID, &t.HelpdeskID, &t.Subject, &t.Status, &t.Priority, &t.RelType, &t.RelID, &t.CreatedAt, &t.UpdatedAt); return &t, err
}

func (r *SupportRepository) GetTicketByID(ctx context.Context, id int64) (*domain.Ticket, error) {
	t, err := scanTicket(r.pool.QueryRow(ctx, "SELECT "+ticketCols+" FROM support_tickets WHERE id = $1", id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) { return nil, appErrors.ErrNotFound }; return t, err
}

func (r *SupportRepository) ListTicketsByClientID(ctx context.Context, cID int64, l, o int) ([]*domain.Ticket, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+ticketCols+" FROM support_tickets WHERE client_id = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", scanTicket, cID, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM support_tickets WHERE client_id = $1", cID), err
}

func (r *SupportRepository) ListTickets(ctx context.Context, l, o int) ([]*domain.Ticket, int, error) {
	res, err := list(ctx, r.pool, "SELECT "+ticketCols+" FROM support_tickets ORDER BY id DESC LIMIT $1 OFFSET $2", scanTicket, l, o)
	return res, total(ctx, r.pool, "SELECT COUNT(*) FROM support_tickets"), err
}

func (r *SupportRepository) CreateTicket(ctx context.Context, t *domain.Ticket, msg *domain.TicketMessage) error {
	tx, err := r.pool.Begin(ctx); if err != nil { return err }; defer tx.Rollback(ctx)
	if t.Status == "" { t.Status = domain.TicketStatusOpen }; if t.Priority == "" { t.Priority = domain.PriorityMedium }
	if err := tx.QueryRow(ctx, `INSERT INTO support_tickets (client_id, helpdesk_id, subject, status, priority, rel_type, rel_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, created_at, updated_at`, t.ClientID, t.HelpdeskID, t.Subject, t.Status, t.Priority, t.RelType, t.RelID).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt); err != nil { return err }
	if msg != nil {
		msg.TicketID = t.ID; _ = tx.QueryRow(ctx, `INSERT INTO support_ticket_messages (ticket_id, admin_id, client_id, content, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) RETURNING id, created_at`, t.ID, msg.AdminID, msg.ClientID, msg.Content, msg.IPAddress).Scan(&msg.ID, &msg.CreatedAt)
	}
	return tx.Commit(ctx)
}

func (r *SupportRepository) AddMessage(ctx context.Context, m *domain.TicketMessage) error {
	return r.pool.QueryRow(ctx, `INSERT INTO support_ticket_messages (ticket_id, admin_id, client_id, content, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) RETURNING id, created_at`, m.TicketID, m.AdminID, m.ClientID, m.Content, m.IPAddress).Scan(&m.ID, &m.CreatedAt)
}

func (r *SupportRepository) UpdateTicketStatus(ctx context.Context, id int64, st domain.TicketStatus) error {
	t, err := r.pool.Exec(ctx, `UPDATE support_tickets SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, st, id)
	if err == nil && t.RowsAffected() == 0 { return appErrors.ErrNotFound }; return err
}

func (r *SupportRepository) GetMessages(ctx context.Context, id int64) ([]*domain.TicketMessage, error) {
	return list(ctx, r.pool, `SELECT id, ticket_id, admin_id, client_id, content, ip_address, created_at FROM support_ticket_messages WHERE ticket_id = $1 ORDER BY id ASC`, func(r pgx.Row) (*domain.TicketMessage, error) {
		var m domain.TicketMessage; err := r.Scan(&m.ID, &m.TicketID, &m.AdminID, &m.ClientID, &m.Content, &m.IPAddress, &m.CreatedAt); return &m, err
	}, id)
}

func (r *SupportRepository) CloseInactiveTickets(ctx context.Context, cut time.Time) (int, error) {
	t, err := r.pool.Exec(ctx, `UPDATE support_tickets SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE status != $1 AND updated_at <= $2`, domain.TicketStatusClosed, cut); return int(t.RowsAffected()), err
}
