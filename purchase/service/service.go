package service

import (
	"context"
	"fmt"
	"time"

	accountingDTO "github.com/aceextension/accounting/dto"
	accountingService "github.com/aceextension/accounting/service"
	"github.com/aceextension/inventory/service"
	"github.com/aceextension/purchase/domain"
	"github.com/aceextension/purchase/dto"
	"github.com/aceextension/purchase/repository"
	"github.com/google/uuid"
)

type PurchaseService interface {
	CreateBill(ctx context.Context, tenantID, userID uuid.UUID, req dto.CreateBillRequest) (*domain.PurchaseBill, error)
	GetBill(ctx context.Context, id uuid.UUID) (*domain.PurchaseBill, error)
	ListBills(ctx context.Context, tenantID uuid.UUID) ([]*domain.PurchaseBill, error)
}

type purchaseService struct {
	repo              repository.PurchaseRepository
	inventoryService  service.InventoryService
	accountingService accountingService.AccountingService
}

func NewPurchaseService(
	repo repository.PurchaseRepository,
	invService service.InventoryService,
	accService accountingService.AccountingService,
) PurchaseService {
	return &purchaseService{
		repo:              repo,
		inventoryService:  invService,
		accountingService: accService,
	}
}

func (s *purchaseService) CreateBill(ctx context.Context, tenantID, userID uuid.UUID, req dto.CreateBillRequest) (*domain.PurchaseBill, error) {
	// 1. Calculate Totals
	var totalAmount float64
	var billItems []domain.PurchaseBillItem

	for _, itemReq := range req.Items {
		total := float64(itemReq.Quantity) * itemReq.UnitPrice
		totalAmount += total

		billItems = append(billItems, domain.PurchaseBillItem{
			ID:        uuid.New(),
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			Total:     total,
			TaxAmount: itemReq.TaxAmount,
		})
	}

	bill := &domain.PurchaseBill{
		ID:          uuid.New(),
		TenantID:    tenantID,
		BillNo:      req.BillNo,
		SupplierID:  req.SupplierID,
		IssueDate:   req.IssueDate,
		DueDate:     req.DueDate,
		TotalAmount: totalAmount,
		Status:      domain.BillStatusReceived, // Assuming direct receipt for now
		Items:       billItems,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 2. Save Bill
	if err := s.repo.CreateBill(ctx, bill); err != nil {
		return nil, err
	}

	// 3. Add Inventory (Stock In)
	// Same issue as Sales - need WarehouseID.
	// Placeholder: uuid.Nil or skip if strict.
	// Commenting out until Warehouse logic defined.
	/*
		warehouseID := uuid.Nil // TODO: Resolve warehouse
		for _, item := range req.Items {
			err := s.inventoryService.AddStock(ctx, tenantID, warehouseID, item.ProductID, float64(item.Quantity), "PURCHASE_BILL", &bill.ID, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to add stock for product %s: %w", item.ProductID, err)
			}
		}
	*/

	// 4. Post Accounting Entry
	// Debit: Inventory Asset or Expense
	// Credit: Accounts Payable (Supplier) or Cash

	journalReq := accountingDTO.CreateJournalEntryRequest{
		FiscalYearID:  uuid.Nil, // Resolve
		Date:          bill.IssueDate,
		Description:   fmt.Sprintf("Purchase Bill #%s", bill.BillNo),
		ReferenceID:   &bill.ID,
		ReferenceType: func() *string { s := "PURCHASE_BILL"; return &s }(),
		Lines: []accountingDTO.JournalLineRequest{
			{
				AccountID:   uuid.Nil, // Expense/Asset Account ID
				Debit:       totalAmount,
				Credit:      0,
				Description: nil,
			},
			{
				AccountID:   req.CreditAccount, // AP or Cash
				Debit:       0,
				Credit:      totalAmount,
				Description: nil,
			},
		},
	}

	// s.accountingService.CreateJournalEntry(ctx, tenantID, userID, journalReq)

	return bill, nil
}

func (s *purchaseService) GetBill(ctx context.Context, id uuid.UUID) (*domain.PurchaseBill, error) {
	return s.repo.GetBill(ctx, id)
}

func (s *purchaseService) ListBills(ctx context.Context, tenantID uuid.UUID) ([]*domain.PurchaseBill, error) {
	return s.repo.ListBills(ctx, tenantID)
}
