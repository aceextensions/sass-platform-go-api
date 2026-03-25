package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
)

type MigrationFile struct {
	FullPath string
	FileName string
}

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

func ensureMigrationTable(ctx context.Context) {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := db.MainPool.Exec(ctx, query)
	if err != nil {
		log.Fatalf("Failed to create schema_migrations table: %v", err)
	}
}

func getAppliedMigrations(ctx context.Context) map[string]bool {
	rows, err := db.MainPool.Query(ctx, "SELECT filename FROM schema_migrations")
	if err != nil {
		log.Fatalf("Failed to fetch applied migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err == nil {
			applied[filename] = true
		}
	}
	return applied
}

// discoverMigrations finds all .sql files in any folder named 'migrations'
func discoverMigrations() []MigrationFile {
	var migrations []MigrationFile
	
	// Determine the workspace root
	root := "."
	if _, err := os.Stat("go.work"); err != nil {
		// If not in root, try one level up (common if running from api/cmd/migrate)
		if _, err := os.Stat("../go.work"); err == nil {
			root = ".."
		} else if _, err := os.Stat("../../go.work"); err == nil {
			root = "../.."
		}
	}

	absRoot, _ := filepath.Abs(root)
	// log.Printf("🔍 Scanning workspace root: %s", absRoot)

	err := filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// Skip hidden directories (like .git, .gemini)
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		// If it's a directory named 'migrations', collect SQL files
		if d.IsDir() && d.Name() == "migrations" {
			// Find all .sql files in this directory
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil
			}

			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
					fullPath := filepath.Join(path, entry.Name())
					// Create a unique key relative to the workspace root
					relPath, _ := filepath.Rel(absRoot, fullPath)
					
					migrations = append(migrations, MigrationFile{
						FullPath: fullPath,
						FileName: relPath,
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Warning: error during migration discovery: %v", err)
	}

	// Sort migrations lexicographically across all modules
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].FileName < migrations[j].FileName
	})

	return migrations
}

func runMigrations(ctx context.Context) {
	fmt.Println("🚀 Discovering all module migrations...")
	ensureMigrationTable(ctx)
	applied := getAppliedMigrations(ctx)
	migrations := discoverMigrations()

	count := 0
	for _, m := range migrations {
		if applied[m.FileName] {
			continue
		}

		fmt.Printf("⏳ Applying %s...\n", m.FileName)
		content, err := os.ReadFile(m.FullPath)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", m.FileName, err)
		}

		tx, err := db.MainPool.Begin(ctx)
		if err != nil {
			log.Fatalf("Failed to start transaction: %v", err)
		}

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("❌ Failed to apply %s: %v", m.FileName, err)
		}

		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", m.FileName); err != nil {
			tx.Rollback(ctx)
			log.Fatalf("❌ Failed to record %s: %v", m.FileName, err)
		}

		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("Failed to commit %s: %v", m.FileName, err)
		}

		fmt.Printf("✅ Applied %s\n", m.FileName)
		count++
	}

	if count == 0 {
		fmt.Println("✨ All module databases are already up to date.")
	} else {
		fmt.Printf("Successfully applied %d migrations across all modules.\n", count)
	}
}

func showStatus(ctx context.Context) {
	ensureMigrationTable(ctx)
	applied := getAppliedMigrations(ctx)
	migrations := discoverMigrations()

	fmt.Println("\nGlobal Migration Status:")
	fmt.Println("------------------------")
	for _, m := range migrations {
		status := "PENDING"
		if applied[m.FileName] {
			status = "APPLIED"
		}
		fmt.Printf("[%s] %s\n", status, m.FileName)
	}
	
	if len(migrations) == 0 {
		fmt.Println("No migration files found in the workspace.")
	}
	fmt.Println("------------------------")
}
