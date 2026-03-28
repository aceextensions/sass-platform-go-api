package db

import (
	"fmt"
	"io/fs"
	"path"

	accountingMig "github.com/aceextension/accounting/migrations"
	auditMig "github.com/aceextension/audit/migrations"
	catalogMig "github.com/aceextension/catalog/migrations"
	coreMig "github.com/aceextension/core/migrations"
	crmMig "github.com/aceextension/crm/migrations"
	fiscalMig "github.com/aceextension/fiscal/migrations"
	inventoryMig "github.com/aceextension/inventory/migrations"
	notificationMig "github.com/aceextension/notification/migrations"
	purchaseMig "github.com/aceextension/purchase/migrations"
	salesMig "github.com/aceextension/sales/migrations"
	subscriptionMig "github.com/aceextension/subscription/migrations"
)

// MigrationProvider is a function type that returns a migration filesystem
type MigrationProvider func() (fs.FS, error)

// GetAllMigrations aggregates all migrations from all modules into a sorted list
func GetAllMigrations() ([]Migration, error) {
	// 1. Define order of modules (critical for dependency management)
	providers := []struct {
		prefix string
		fs     fs.FS
	}{
		{"core", coreMig.FS},           // Base tables (tenants, users)
		{"audit", auditMig.FS},         // Monitoring
		{"fiscal", fiscalMig.FS},       // Timeframes
		{"catalog", catalogMig.FS},     // Products
		{"inventory", inventoryMig.FS}, // Stock (depends on products)
		{"crm", crmMig.FS},             // Customers/Suppliers
		{"notification", notificationMig.FS},
		{"subscription", subscriptionMig.FS},
		{"accounting", accountingMig.FS},
		{"sales", salesMig.FS},         // Transactions (depends on customers/products)
		{"purchase", purchaseMig.FS},   // Transactions (depends on suppliers/products)
	}

	var all []Migration
	for _, p := range providers {
		err := fs.WalkDir(p.fs, ".", func(filePath string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || path.Ext(filePath) != ".sql" {
				return nil
			}

			content, err := fs.ReadFile(p.fs, filePath)
			if err != nil {
				return fmt.Errorf("failed to read %s: %w", filePath, err)
			}

			all = append(all, Migration{
				ID:  fmt.Sprintf("%s/%s", p.prefix, filePath),
				SQL: string(content),
				// For now, version is implicitly handled by filenames within modules,
				// and by the provider order in this list.
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return all, nil
}
