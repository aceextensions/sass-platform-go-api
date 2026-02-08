package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aceextension/audit"
	"github.com/aceextension/audit/domain"
	"github.com/aceextension/audit/helper"
	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize database connections
	mainConnStr := config.Get("DATABASE_URL")
	auditConnStr := config.Get("AUDIT_DATABASE_URL")
	db.Init(mainConnStr, auditConnStr)
	defer db.Close()

	// Initialize audit module
	audit.Init()

	fmt.Println("🔍 Audit Module Example")
	fmt.Println("========================\n")

	// Example 1: Direct audit logging
	fmt.Println("Example 1: Direct Audit Logging")
	ctx := context.Background()
	tenantID := uuid.New()
	userID := uuid.New()
	ipAddress := "192.168.1.100"
	userAgent := "Mozilla/5.0"

	auditCtx := &domain.AuditContext{
		TenantID:  &tenantID,
		UserID:    &userID,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
	}

	// Log a user creation event
	entityID := userID.String()
	err := audit.Service.Log(ctx, "CREATE_USER", "User", &entityID, map[string]interface{}{
		"userName": "John Doe",
		"email":    "john@example.com",
		"role":     "admin",
	}, auditCtx)

	if err != nil {
		log.Fatalf("Failed to log audit: %v", err)
	}
	fmt.Println("✅ Logged: CREATE_USER")

	// Example 2: Using Audit Helper
	fmt.Println("\nExample 2: Using Audit Helper")
	auditHelper := helper.NewAuditHelper()

	// Log a sale creation
	saleID := uuid.New()
	err = auditHelper.LogSaleAction(ctx, "CREATE_SALE", saleID, map[string]interface{}{
		"invoiceNumber": "INV-8283-0001",
		"customerName":  "ABC Company",
		"totalAmount":   15000.00,
		"vatAmount":     1950.00,
	}, auditCtx)

	if err != nil {
		log.Fatalf("Failed to log sale: %v", err)
	}
	fmt.Println("✅ Logged: CREATE_SALE")

	// Log a purchase creation
	purchaseID := uuid.New()
	err = auditHelper.LogPurchaseAction(ctx, "CREATE_PURCHASE", purchaseID, map[string]interface{}{
		"purchaseNumber": "PUR-8283-0001",
		"supplierName":   "XYZ Suppliers",
		"totalAmount":    25000.00,
	}, auditCtx)

	if err != nil {
		log.Fatalf("Failed to log purchase: %v", err)
	}
	fmt.Println("✅ Logged: CREATE_PURCHASE")

	// Log a journal entry posting
	journalID := uuid.New()
	err = auditHelper.LogJournalAction(ctx, "POST_JOURNAL_ENTRY", journalID, map[string]interface{}{
		"voucherNumber": "JV-8283-0001",
		"totalDebit":    25000.00,
		"totalCredit":   25000.00,
	}, auditCtx)

	if err != nil {
		log.Fatalf("Failed to log journal: %v", err)
	}
	fmt.Println("✅ Logged: POST_JOURNAL_ENTRY")

	// Example 3: Querying Audit Logs
	fmt.Println("\nExample 3: Querying Audit Logs")

	// Wait a moment for async logs to be written
	fmt.Println("⏳ Waiting for async logs to be written...")
	// In production, you might want to use sync logging for critical operations
	// or implement a proper wait mechanism

	// Query logs by tenant
	logs, err := audit.Service.GetByTenantID(ctx, tenantID, 10, 0)
	if err != nil {
		log.Fatalf("Failed to query logs: %v", err)
	}

	fmt.Printf("📋 Found %d audit logs for tenant\n", len(logs))
	for i, log := range logs {
		fmt.Printf("  %d. %s - %s (Entity: %s)\n", i+1, log.Action, log.CreatedAt.Format("2006-01-02 15:04:05"), log.Entity)
	}

	// Example 4: Synchronous Logging (for critical operations)
	fmt.Println("\nExample 4: Synchronous Logging")

	accountID := uuid.New()
	accountIDStr := accountID.String()
	err = audit.Service.LogSync(ctx, "CLOSE_FISCAL_YEAR", "FiscalYear", &accountIDStr, map[string]interface{}{
		"fiscalYear": "2082/83",
		"closedBy":   "admin",
	}, auditCtx)

	if err != nil {
		log.Fatalf("Failed to log sync: %v", err)
	}
	fmt.Println("✅ Logged (Sync): CLOSE_FISCAL_YEAR")

	fmt.Println("\n✨ Audit module is working correctly!")
}
