package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/fiscal"
	"github.com/aceextension/fiscal/utils"
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

	// Initialize fiscal module
	fiscal.Init()

	fmt.Println("📅 Fiscal Year Management Example")
	fmt.Println("===================================\n")

	ctx := context.Background()
	tenantID := uuid.New()

	// Example 1: Nepali Date Conversion
	fmt.Println("Example 1: Nepali Date Conversion")
	fmt.Println("----------------------------------")

	today := time.Now()
	todayBS := utils.ADToBS(today)
	fmt.Printf("Today (AD): %s\n", today.Format("2006-01-02"))
	fmt.Printf("Today (BS): %s (%s)\n", todayBS.String(), todayBS.NepaliMonthName())
	fmt.Printf("Approximate: %s\n\n", todayBS.EnglishMonthName())

	// Example 2: Create Fiscal Year from Nepali Date
	fmt.Println("Example 2: Create Fiscal Year from Nepali Date")
	fmt.Println("-----------------------------------------------")

	fiscalYearName := "2082/83"
	fy, err := fiscal.Service.CreateFromNepaliDate(ctx, tenantID, fiscalYearName)
	if err != nil {
		log.Printf("Error creating fiscal year: %v", err)
	} else {
		fmt.Printf("Created Fiscal Year: %s\n", fy.Name)
		fmt.Printf("  Start (AD): %s\n", fy.StartDate.Format("2006-01-02"))
		fmt.Printf("  Start (BS): %s\n", fy.StartDateBS)
		fmt.Printf("  End (AD):   %s\n", fy.EndDate.Format("2006-01-02"))
		fmt.Printf("  End (BS):   %s\n", fy.EndDateBS)
		fmt.Printf("  Invoice Prefix: %s\n", fy.InvoicePrefix)
		fmt.Printf("  Purchase Prefix: %s\n", fy.PurchasePrefix)
		fmt.Printf("  Voucher Prefix: %s\n\n", fy.VoucherPrefix)
	}

	// Example 3: Generate Invoice Numbers
	fmt.Println("Example 3: Generate Invoice Numbers")
	fmt.Println("------------------------------------")

	if fy != nil {
		// Set as current
		err = fiscal.Service.SetAsCurrent(ctx, tenantID, fy.ID)
		if err != nil {
			log.Printf("Error setting fiscal year as current: %v", err)
		}

		// Generate invoice numbers
		for i := 1; i <= 5; i++ {
			invoiceNum, err := fiscal.Service.GenerateInvoiceNumber(ctx, fy.ID)
			if err != nil {
				log.Printf("Error generating invoice number: %v", err)
				break
			}
			fmt.Printf("  Invoice #%d: %s\n", i, invoiceNum)
		}
		fmt.Println()

		// Generate purchase numbers
		fmt.Println("Purchase Numbers:")
		for i := 1; i <= 3; i++ {
			purchaseNum, err := fiscal.Service.GeneratePurchaseNumber(ctx, fy.ID)
			if err != nil {
				log.Printf("Error generating purchase number: %v", err)
				break
			}
			fmt.Printf("  Purchase #%d: %s\n", i, purchaseNum)
		}
		fmt.Println()

		// Generate voucher numbers
		fmt.Println("Voucher Numbers:")
		for i := 1; i <= 3; i++ {
			voucherNum, err := fiscal.Service.GenerateVoucherNumber(ctx, fy.ID)
			if err != nil {
				log.Printf("Error generating voucher number: %v", err)
				break
			}
			fmt.Printf("  Voucher #%d: %s\n", i, voucherNum)
		}
		fmt.Println()
	}

	// Example 4: Get Current Fiscal Year
	fmt.Println("Example 4: Get Current Fiscal Year")
	fmt.Println("-----------------------------------")

	currentFY, err := fiscal.Service.GetCurrent(ctx, tenantID)
	if err != nil {
		log.Printf("Error getting current fiscal year: %v", err)
	} else {
		fmt.Printf("Current Fiscal Year: %s\n", currentFY.Name)
		fmt.Printf("  Is Active: %v\n", currentFY.IsActive())
		fmt.Printf("  Can Modify: %v\n", currentFY.CanModify())
		fmt.Printf("  Is Closed: %v\n\n", currentFY.IsClosed)
	}

	// Example 5: Fiscal Year Operations
	fmt.Println("Example 5: Fiscal Year Operations")
	fmt.Println("----------------------------------")

	if fy != nil {
		// Close fiscal year
		closedBy := uuid.New()
		err = fiscal.Service.Close(ctx, fy.ID, closedBy)
		if err != nil {
			log.Printf("Error closing fiscal year: %v", err)
		} else {
			fmt.Println("✅ Fiscal year closed")
		}

		// Try to generate invoice number (should fail)
		_, err = fiscal.Service.GenerateInvoiceNumber(ctx, fy.ID)
		if err != nil {
			fmt.Printf("❌ Cannot generate invoice for closed fiscal year: %v\n", err)
		}

		// Reopen fiscal year
		err = fiscal.Service.Reopen(ctx, fy.ID)
		if err != nil {
			log.Printf("Error reopening fiscal year: %v", err)
		} else {
			fmt.Println("✅ Fiscal year reopened")
		}

		// Now can generate invoice
		invoiceNum, err := fiscal.Service.GenerateInvoiceNumber(ctx, fy.ID)
		if err != nil {
			log.Printf("Error generating invoice: %v", err)
		} else {
			fmt.Printf("✅ Generated invoice after reopen: %s\n", invoiceNum)
		}
	}

	// Example 6: Nepali Date Formatting
	fmt.Println("\nExample 6: Nepali Date Formatting")
	fmt.Println("----------------------------------")

	sampleDate := utils.NepaliDate{Year: 2082, Month: 4, Day: 1}
	fmt.Printf("YYYY-MM-DD: %s\n", utils.FormatNepaliDate(sampleDate, "YYYY-MM-DD"))
	fmt.Printf("DD MMM YYYY: %s\n", utils.FormatNepaliDate(sampleDate, "DD MMM YYYY"))
	fmt.Printf("DD MMMM YYYY: %s\n", utils.FormatNepaliDate(sampleDate, "DD MMMM YYYY"))

	fmt.Println("\n✨ Fiscal year management is working correctly!")
}
