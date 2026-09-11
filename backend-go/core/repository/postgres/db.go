package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL); if err != nil { return nil, fmt.Errorf("unable to parse database config: %w", err) }
	cfg.MaxConns, cfg.MinConns = 25, 5
	cfg.MaxConnLifetime, cfg.MaxConnIdleTime = time.Hour, 30*time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg); if err != nil { return nil, fmt.Errorf("unable to create pool: %w", err) }
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second); defer cancel()
	if err := pool.Ping(pingCtx); err != nil { pool.Close(); return nil, fmt.Errorf("unable to ping database: %w", err) }
	log.Println("✅ Successfully connected to PostgreSQL database pool.")
	return pool, nil
}

// Helper: list scans multiple rows into a slice using a scanner function
func list[T any](ctx context.Context, p *pgxpool.Pool, q string, sc func(pgx.Row) (*T, error), args ...any) ([]*T, error) {
	rows, err := p.Query(ctx, q, args...); if err != nil { return nil, err }; defer rows.Close()
	var res []*T
	for rows.Next() {
		it, err := sc(rows); if err != nil { return nil, err }
		res = append(res, it)
	}
	return res, nil
}

// Helper: total returns the row count for a given query
func total(ctx context.Context, p *pgxpool.Pool, q string, args ...any) int {
	var count int
	err := p.QueryRow(ctx, q, args...).Scan(&count)
	if err != nil {
		log.Printf("⚠️  Database count query failed: %v (Query: %s)", err, q)
		return 0
	}
	return count
}
