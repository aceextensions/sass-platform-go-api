package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
)

func main() {
	// 1. Load Config & DB
	conf := config.Load()
	db.Init(conf.DatabaseURL, conf.AuditDatabaseURL)
	defer db.Close()

	ctx := context.Background()

	// 2. Determine Action
	action := "up"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	switch action {
	case "up":
		runMigrations(ctx)
	case "status":
		showStatus(ctx)
	default:
		fmt.Printf("Unknown command: %s\nUsage: go run ./api/cmd/migrate/main.go [up|status]\n", action)
		os.Exit(1)
	}
}

func runMigrations(ctx context.Context) {
	fmt.Println("🚀 Running unified auto-migrations...")
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}
	fmt.Println("✨ All module databases are now up to date.")
}

func showStatus(ctx context.Context) {
	if db.MainPool == nil {
		log.Fatal("Database pool not initialized")
	}

	migrator := db.NewMigrator(db.MainPool)
	if err := migrator.EnsureSchemaTable(ctx); err != nil {
		log.Fatalf("Failed to ensure schema table: %v", err)
	}

	applied, err := migrator.GetAppliedMigrations(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch applied migrations: %v", err)
	}

	migrations, err := db.GetAllMigrations()
	if err != nil {
		log.Fatalf("Failed to gather migrations: %v", err)
	}

	fmt.Println("\nUnified Migration Status:")
	fmt.Println("------------------------")
	for _, m := range migrations {
		status := "PENDING"
		if applied[m.ID] {
			status = "APPLIED"
		}
		fmt.Printf("[%s] %s\n", status, m.ID)
	}
	
	if len(migrations) == 0 {
		fmt.Println("No migration files found in the workspace.")
	}
	fmt.Println("------------------------")
}
