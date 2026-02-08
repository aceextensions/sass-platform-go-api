package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/crm"
	"github.com/aceextension/crm/domain"
	"github.com/aceextension/fiscal"
	"github.com/google/uuid"
)

func main() {
	// Initialize configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize modules
	fiscal.Init()
	crm.Init()

	ctx := context.Background()
	tenantID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000") // Example tenant ID

	fmt.Println("=== CRM Module Example ===\n")

	// Example 1: Create Customer
	fmt.Println("1. Creating Customer...")
	customer := domain.NewCustomer(tenantID, "ABC Trading Company")
	email := "contact@abctrading.com"
	phone := "9841234567"
	customer.Email = &email
	customer.Phone = &phone
	customer.CustomerType = domain.CustomerTypeBusiness

	// Add custom attributes
	customer.SetPANNumber("123456789")
	customer.SetVATNumber("987654321")
	customer.SetCreditLimit(100000.00)
	customer.SetAddress("Thamel, Kathmandu, Nepal")
	customer.SetCustomAttribute("loyalty_tier", "gold")
	customer.SetCustomAttribute("payment_terms", "30_days")

	if err := crm.CustomerService.Create(ctx, customer); err != nil {
		log.Fatalf("Failed to create customer: %v", err)
	}

	fmt.Printf("✓ Customer created: %s (%s)\n", customer.Name, customer.CustomerCode)
	fmt.Printf("  PAN: %s\n", customer.GetPANNumber())
	fmt.Printf("  Credit Limit: %.2f\n", customer.GetCreditLimit())
	fmt.Printf("  Loyalty Tier: %s\n\n", customer.GetCustomString("loyalty_tier"))

	// Example 2: Create Supplier
	fmt.Println("2. Creating Supplier...")
	supplier := domain.NewSupplier(tenantID, "XYZ Suppliers Pvt. Ltd.")
	supplierEmail := "info@xyzsuppliers.com"
	supplierPhone := "9851234567"
	supplier.Email = &supplierEmail
	supplier.Phone = &supplierPhone
	supplier.SupplierType = domain.SupplierTypeLocal

	// Add custom attributes
	supplier.SetPANNumber("987654321")
	supplier.SetPaymentTerms("15_days")
	supplier.SetLeadTimeDays(7)
	supplier.SetMinimumOrderValue(5000.00)
	supplier.SetAddress("Patan, Lalitpur, Nepal")
	supplier.SetCustomAttribute("bank_name", "Nepal Bank Limited")
	supplier.SetCustomAttribute("bank_account", "1234567890")

	if err := crm.SupplierService.Create(ctx, supplier); err != nil {
		log.Fatalf("Failed to create supplier: %v", err)
	}

	fmt.Printf("✓ Supplier created: %s (%s)\n", supplier.Name, supplier.SupplierCode)
	fmt.Printf("  Payment Terms: %s\n", supplier.GetPaymentTerms())
	fmt.Printf("  Lead Time: %d days\n", supplier.GetLeadTimeDays())
	fmt.Printf("  Min Order Value: %.2f\n\n", supplier.GetMinimumOrderValue())

	// Example 3: Search Customers
	fmt.Println("3. Searching Customers...")
	customers, err := crm.CustomerService.Search(ctx, tenantID, "ABC", 10, 0)
	if err != nil {
		log.Fatalf("Failed to search customers: %v", err)
	}

	fmt.Printf("✓ Found %d customer(s) matching 'ABC'\n", len(customers))
	for _, c := range customers {
		fmt.Printf("  - %s (%s)\n", c.Name, c.CustomerCode)
	}
	fmt.Println()

	// Example 4: Get Customer by Code
	fmt.Println("4. Getting Customer by Code...")
	retrievedCustomer, err := crm.CustomerService.GetByCode(ctx, tenantID, customer.CustomerCode)
	if err != nil {
		log.Fatalf("Failed to get customer: %v", err)
	}

	fmt.Printf("✓ Retrieved customer: %s\n", retrievedCustomer.Name)
	fmt.Printf("  Email: %s\n", *retrievedCustomer.Email)
	fmt.Printf("  Status: %s\n", retrievedCustomer.Status)
	fmt.Printf("  Custom Attributes: %d\n\n", len(retrievedCustomer.CustomAttributes))

	// Example 5: Update Customer
	fmt.Println("5. Updating Customer...")
	retrievedCustomer.SetCreditLimit(150000.00)
	retrievedCustomer.SetCustomAttribute("loyalty_tier", "platinum")

	if err := crm.CustomerService.Update(ctx, retrievedCustomer); err != nil {
		log.Fatalf("Failed to update customer: %v", err)
	}

	fmt.Printf("✓ Customer updated\n")
	fmt.Printf("  New Credit Limit: %.2f\n", retrievedCustomer.GetCreditLimit())
	fmt.Printf("  New Loyalty Tier: %s\n\n", retrievedCustomer.GetCustomString("loyalty_tier"))

	// Example 6: Count Customers and Suppliers
	fmt.Println("6. Counting Entities...")
	customerCount, err := crm.CustomerService.Count(ctx, tenantID)
	if err != nil {
		log.Fatalf("Failed to count customers: %v", err)
	}

	supplierCount, err := crm.SupplierService.Count(ctx, tenantID)
	if err != nil {
		log.Fatalf("Failed to count suppliers: %v", err)
	}

	fmt.Printf("✓ Total Customers: %d\n", customerCount)
	fmt.Printf("✓ Total Suppliers: %d\n\n", supplierCount)

	// Example 7: List All Customers
	fmt.Println("7. Listing All Customers...")
	allCustomers, err := crm.CustomerService.GetByTenantID(ctx, tenantID, 10, 0)
	if err != nil {
		log.Fatalf("Failed to list customers: %v", err)
	}

	fmt.Printf("✓ Customers for tenant:\n")
	for _, c := range allCustomers {
		fmt.Printf("  - %s (%s) - %s\n", c.Name, c.CustomerCode, c.Status)
	}

	fmt.Println("\n=== CRM Module Example Complete ===")
}
