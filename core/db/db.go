package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// QueryExecutor is an interface that matches both *pgxpool.Pool and pgx.Tx
type QueryExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

var (
	MainPool  *pgxpool.Pool
	AuditPool *pgxpool.Pool
)

func Init(mainConnStr, auditConnStr string) {
	var err error
	ctx := context.Background()

	// Initialize Main Database Pool
	MainPool, err = pgxpool.New(ctx, mainConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to main database: %v\n", err)
	}

	// Verify connection
	if err := MainPool.Ping(ctx); err != nil {
		log.Fatalf("Main database ping failed: %v\n", err)
	}
	fmt.Println("🚀 Connected to Main Database")

	// Initialize Audit Database Pool
	AuditPool, err = pgxpool.New(ctx, auditConnStr)
	if err != nil {
		log.Fatalf("Unable to connect to audit database: %v\n", err)
	}

	// Verify connection
	if err := AuditPool.Ping(ctx); err != nil {
		log.Fatalf("Audit database ping failed: %v\n", err)
	}
	fmt.Println("🚀 Connected to Audit Database")
}

func Close() {
	if MainPool != nil {
		MainPool.Close()
	}
	if AuditPool != nil {
		AuditPool.Close()
	}
}

func BeginFunc(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := MainPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = fn(tx)
	return err
}
