package importer

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/damarkuncoro/FOSSBilling/backend-go/core/domain"
	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/decimal"
)

type LegacyImporter struct {
	clientRepo  domain.ClientRepository
	orderRepo   domain.OrderRepository
	invoiceRepo domain.InvoiceRepository
}

func NewLegacyImporter(cr domain.ClientRepository, or domain.OrderRepository, ir domain.InvoiceRepository) *LegacyImporter {
	return &LegacyImporter{clientRepo: cr, orderRepo: or, invoiceRepo: ir}
}

func (i *LegacyImporter) ImportClients(ctx context.Context, sourceDB *sql.DB) (int, error) {
	rows, err := sourceDB.QueryContext(ctx, "SELECT id, email, first_name, last_name, currency, created_at FROM client")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var c domain.Client
		var createdAt string
		if err := rows.Scan(&c.ID, &c.Email, &c.FirstName, &c.LastName, &c.Currency, &createdAt); err != nil {
			continue
		}

		c.Status = domain.ClientStatusActive
		c.PasswordHash = "migrated_account_please_reset"

		if err := i.clientRepo.Create(ctx, &c); err != nil {
			continue
		}
		count++
	}
	return count, nil
}

func (i *LegacyImporter) ImportInvoices(ctx context.Context, sourceDB *sql.DB) (int, error) {
	rows, err := sourceDB.QueryContext(ctx, "SELECT id, client_id, serie, nr, currency, subtotal, tax, total, status FROM invoice")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var inv domain.Invoice
		var sub, tx, tot float64
		if err := rows.Scan(&inv.ID, &inv.ClientID, &inv.Serie, &inv.Nr, &inv.Currency, &sub, &tx, &tot, &inv.Status); err != nil {
			continue
		}

		inv.Subtotal = decimal.FromFloat(sub)
		inv.Tax = decimal.FromFloat(tx)
		inv.Total = decimal.FromFloat(tot)

		if err := i.invoiceRepo.Create(ctx, &inv, []domain.InvoiceItem{}); err == nil {
			count++
		}
	}
	return count, nil
}

func (i *LegacyImporter) RunFullMigration(ctx context.Context, sourceDSN string) error {
	sdb, err := sql.Open("mysql", sourceDSN)
	if err != nil {
		return err
	}
	defer sdb.Close()

	log.Println("🏁 Starting Legacy Migration (PHP -> Go)...")

	cCount, _ := i.ImportClients(ctx, sdb)
	log.Printf("✅ Imported %d Clients", cCount)

	iCount, _ := i.ImportInvoices(ctx, sdb)
	log.Printf("✅ Imported %d Invoices", iCount)

	log.Println("✨ Migration Completed!")
	return nil
}
