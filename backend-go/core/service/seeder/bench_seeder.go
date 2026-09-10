package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/damarkuncoro/FOSSBilling/backend-go/pkg/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BenchSeeder struct {
	pool *pgxpool.Pool
}

func NewBenchSeeder(pool *pgxpool.Pool) *BenchSeeder {
	return &BenchSeeder{pool: pool}
}

func (s *BenchSeeder) SeedLargeScale(ctx context.Context, clientCount, invoiceCount int) error {
	log.Printf("🚀 Starting large-scale seeding: %d clients, %d invoices...", clientCount, invoiceCount)
	start := time.Now()

	// 1. Seed Clients in batches
	clientIDs, err := s.seedClients(ctx, clientCount)
	if err != nil {
		return fmt.Errorf("failed to seed clients: %w", err)
	}

	// 2. Seed Invoices in batches
	if err := s.seedInvoices(ctx, invoiceIDs(clientIDs, invoiceCount)); err != nil {
		return fmt.Errorf("failed to seed invoices: %w", err)
	}

	log.Printf("✅ Seeding completed in %v", time.Since(start))
	return nil
}

func (s *BenchSeeder) seedClients(ctx context.Context, count int) ([]int64, error) {
	batchSize := 1000
	var ids []int64
	pwdHash, _ := auth.HashPassword("Password123!")

	for i := 0; i < count; i += batchSize {
		currentBatch := batchSize
		if i+batchSize > count {
			currentBatch = count - i
		}

		batch := &pgx.Batch{}
		for j := 0; j < currentBatch; j++ {
			idx := i + j
			email := fmt.Sprintf("bench_client_%d@example.com", idx)
			batch.Queue(`INSERT INTO clients (email, password_hash, first_name, last_name, status, currency)
				VALUES ($1, $2, $3, $4, 'active', 'USD') RETURNING id`,
				email, pwdHash, fmt.Sprintf("Bench%d", idx), "User")
		}

		br := s.pool.SendBatch(ctx, batch)
		for j := 0; j < currentBatch; j++ {
			var id int64
			if err := br.QueryRow().Scan(&id); err != nil {
				br.Close()
				return nil, err
			}
			ids = append(ids, id)
		}
		br.Close()
		log.Printf("   ...seeded %d/%d clients", len(ids), count)
	}
	return ids, nil
}

func (s *BenchSeeder) seedInvoices(ctx context.Context, clientInvoices []clientInvoiceMap) error {
	batchSize := 1000
	count := len(clientInvoices)

	for i := 0; i < count; i += batchSize {
		currentBatch := batchSize
		if i+batchSize > count {
			currentBatch = count - i
		}

		batch := &pgx.Batch{}
		for j := 0; j < currentBatch; j++ {
			m := clientInvoices[i+j]
			nr := fmt.Sprintf("%d%d", time.Now().Unix(), i+j)
			batch.Queue(`INSERT INTO invoices (serie, nr, client_id, status, subtotal, tax, total, due_at)
				VALUES ('BCH', $1, $2, 'unpaid', 10000, 0, 10000, CURRENT_DATE + INTERVAL '14 days')`,
				nr, m.ClientID)
		}

		br := s.pool.SendBatch(ctx, batch)
		if err := br.Close(); err != nil {
			return err
		}
		log.Printf("   ...seeded %d/%d invoices", i+currentBatch, count)
	}
	return nil
}

type clientInvoiceMap struct {
	ClientID int64
}

func invoiceIDs(clientIDs []int64, totalInvoices int) []clientInvoiceMap {
	res := make([]clientInvoiceMap, totalInvoices)
	for i := 0; i < totalInvoices; i++ {
		res[i] = clientInvoiceMap{
			ClientID: clientIDs[rand.Intn(len(clientIDs))],
		}
	}
	return res
}
