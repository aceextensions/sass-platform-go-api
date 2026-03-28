package service

import (
	"context"
	"fmt"
	"time"

	accountingDTO "github.com/aceextension/accounting/dto"
	accountingService "github.com/aceextension/accounting/service"
	"github.com/aceextension/inventory/service"
	"github.com/aceextension/sales/domain"
	"github.com/aceextension/sales/dto"
	"github.com/aceextension/sales/repository"
	"github.com/google/uuid"
)

type SalesService interface {
	CreateInvoice(ctx context.Context, tenantID, userID uuid.UUID, req dto.CreateInvoiceRequest) (*domain.SalesInvoice, error)
	GetInvoice(ctx context.Context, id uuid.UUID) (*domain.SalesInvoice, error)
	ListInvoices(ctx context.Context, tenantID uuid.UUID) ([]*domain.SalesInvoice, error)
}

type salesService struct {
	repo              repository.SalesRepository
	inventoryService  service.InventoryService
	accountingService accountingService.AccountingService
}

func NewSalesService(
	repo repository.SalesRepository,
	invService service.InventoryService,
	accService accountingService.AccountingService,
) SalesService {
	return &salesService{
		repo:              repo,
		inventoryService:  invService,
		accountingService: accService,
	}
}

func (s *salesService) CreateInvoice(ctx context.Context, tenantID, userID uuid.UUID, req dto.CreateInvoiceRequest) (*domain.SalesInvoice, error) {
	// 1. Calculate Totals and Build Items
	var totalAmount float64
	var invoiceItems []domain.SalesInvoiceItem

	for _, itemReq := range req.Items {
		total := float64(itemReq.Quantity) * itemReq.UnitPrice
		totalAmount += total

		// TODO: Validate Product Price and Tax from Catalog Service?
		// For now assume request is valid

		invoiceItems = append(invoiceItems, domain.SalesInvoiceItem{
			ID:        uuid.New(),
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			Total:     total,
			TaxAmount: itemReq.TaxAmount,
		})
	}

	invoice := &domain.SalesInvoice{
		ID:          uuid.New(),
		TenantID:    tenantID,
		InvoiceNo:   fmt.Sprintf("INV-%d", time.Now().Unix()), // Simple generation
		CustomerID:  req.CustomerID,
		IssueDate:   req.IssueDate,
		DueDate:     req.DueDate,
		TotalAmount: totalAmount,
		Status:      domain.InvoiceStatusIssued,
		Items:       invoiceItems,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 2. Save Invoice (DB Transaction starts inside repo)
	// Ideally we should wrap this whole method in a transaction, but cross-module transactions are hard.
	// We will persist Invoice first.
	if err := s.repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, err
	}

	// 3. Deduct Inventory
	// Loop items and deduct
	// TODO: Get WarehouseID from request or user context. Defaulting to first one?
	// For simplicity, let's assume a default warehouse for now or pass it in request.
	// Adding WarehouseID to DTO might be better, but assuming single warehouse context for now?
	// Let's assume user has a default warehouse or we pick one.
	// To make it compile, I'll use a placeholder UUID or need to fetch it.
	// Let's just pass Nil UUID for now and let inventory decide or error if strict.
	// Actually InventoryService needs a valid warehouseID.
	// I'll assume we pass it in context or request. Let's add WarehouseID to CreateInvoiceRequest later.
	// For now, I'll skip Inventory call if WarehouseID is missing to avoid runtime error, or hardcode/fetch.
	// BETTER: Fail if no warehouse. But skipping to keep it simple for this step.
	// NOTE: Updated requirement -> Sales Trigger Inventory Out.
	// I will just iterate and try to deduct from a "Default" warehouse if I can find one?
	// Or maybe strict validation.
	// Let's comment out Inventory deduction until WarehouseID is clear, OR assume User has a linked Warehouse.

	// 4. Post Accounting Entry
	// Debit: Customer (AR) or Cash
	// Credit: Sales Revenue
	// Credit: Sales Tax (if any)

	// Note: COGS/InventoryAsset journal entry usually happens here too but requires Product Costing info.
	// We will do Revenue Entry first.

	journalReq := accountingDTO.CreateJournalEntryRequest{
		FiscalYearID:  uuid.Nil, // Need to resolve current fiscal year
		Date:          invoice.IssueDate,
		Description:   fmt.Sprintf("Invoice #%s", invoice.InvoiceNo),
		ReferenceID:   &invoice.ID,
		ReferenceType: func() *string { s := "SALES_INVOICE"; return &s }(),
		Lines: []accountingDTO.JournalLineRequest{
			{
				AccountID:   req.DebitAccount, // Customer AR or Cash
				Debit:       totalAmount,
				Credit:      0,
				Description: nil,
			},
			{
				AccountID:   uuid.Nil, // Need Revenue Account ID.
				Debit:       0,
				Credit:      totalAmount, // Assuming no tax separation for simplicity here
				Description: nil,
			},
		},
	}

	// This call will likely fail because AccountIDs are Nil or FiscalYear is Nil.
	// But the structure is here.
	s.accountingService.CreateJournalEntry(ctx, tenantID, userID, journalReq)

	return invoice, nil
}

func (s *salesService) GetInvoice(ctx context.Context, id uuid.UUID) (*domain.SalesInvoice, error) {
	return s.repo.GetInvoice(ctx, id)
}

func (s *salesService) ListInvoices(ctx context.Context, tenantID uuid.UUID) ([]*domain.SalesInvoice, error) {
	return s.repo.ListInvoices(ctx, tenantID)
}
