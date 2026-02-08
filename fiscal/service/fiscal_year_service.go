package service

import (
	"context"
	"fmt"
	"time"

	"github.com/aceextension/fiscal/domain"
	"github.com/aceextension/fiscal/repository"
	"github.com/aceextension/fiscal/utils"
	"github.com/google/uuid"
)

// FiscalYearService defines the interface for fiscal year operations
type FiscalYearService interface {
	// Create creates a new fiscal year
	Create(ctx context.Context, tenantID uuid.UUID, name string, startDate, endDate time.Time) (*domain.FiscalYear, error)

	// CreateFromNepaliDate creates a fiscal year from Nepali date
	CreateFromNepaliDate(ctx context.Context, tenantID uuid.UUID, fiscalYearName string) (*domain.FiscalYear, error)

	// GetByID retrieves a fiscal year by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FiscalYear, error)

	// GetByTenantID retrieves all fiscal years for a tenant
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.FiscalYear, error)

	// GetCurrent retrieves the current fiscal year for a tenant
	GetCurrent(ctx context.Context, tenantID uuid.UUID) (*domain.FiscalYear, error)

	// SetAsCurrent sets a fiscal year as current
	SetAsCurrent(ctx context.Context, tenantID, fiscalYearID uuid.UUID) error

	// Close closes a fiscal year
	Close(ctx context.Context, fiscalYearID, closedBy uuid.UUID) error

	// Reopen reopens a closed fiscal year
	Reopen(ctx context.Context, fiscalYearID uuid.UUID) error

	// Delete deletes a fiscal year
	Delete(ctx context.Context, fiscalYearID uuid.UUID) error

	// GenerateInvoiceNumber generates the next invoice number
	GenerateInvoiceNumber(ctx context.Context, fiscalYearID uuid.UUID) (string, error)

	// GeneratePurchaseNumber generates the next purchase number
	GeneratePurchaseNumber(ctx context.Context, fiscalYearID uuid.UUID) (string, error)

	// GenerateVoucherNumber generates the next voucher number
	GenerateVoucherNumber(ctx context.Context, fiscalYearID uuid.UUID) (string, error)
}

// fiscalYearService implements FiscalYearService
type fiscalYearService struct {
	repo repository.FiscalYearRepository
}

// NewFiscalYearService creates a new fiscal year service
func NewFiscalYearService(repo repository.FiscalYearRepository) FiscalYearService {
	return &fiscalYearService{
		repo: repo,
	}
}

// Create creates a new fiscal year
func (s *fiscalYearService) Create(ctx context.Context, tenantID uuid.UUID, name string, startDate, endDate time.Time) (*domain.FiscalYear, error) {
	// Convert dates to Nepali
	startBS := utils.ADToBS(startDate)
	endBS := utils.ADToBS(endDate)

	// Create fiscal year
	fy := domain.NewFiscalYear(tenantID, name, startDate, endDate, startBS.String(), endBS.String())

	// Save to database
	if err := s.repo.Create(ctx, fy); err != nil {
		return nil, fmt.Errorf("failed to create fiscal year: %w", err)
	}

	return fy, nil
}

// CreateFromNepaliDate creates a fiscal year from Nepali fiscal year name
func (s *fiscalYearService) CreateFromNepaliDate(ctx context.Context, tenantID uuid.UUID, fiscalYearName string) (*domain.FiscalYear, error) {
	// Get fiscal year dates
	startBS, endBS, startAD, endAD := utils.GetFiscalYearDates(fiscalYearName)

	// Create fiscal year
	fy := domain.NewFiscalYear(tenantID, fiscalYearName, startAD, endAD, startBS.String(), endBS.String())

	// Save to database
	if err := s.repo.Create(ctx, fy); err != nil {
		return nil, fmt.Errorf("failed to create fiscal year: %w", err)
	}

	return fy, nil
}

// GetByID retrieves a fiscal year by ID
func (s *fiscalYearService) GetByID(ctx context.Context, id uuid.UUID) (*domain.FiscalYear, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByTenantID retrieves all fiscal years for a tenant
func (s *fiscalYearService) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.FiscalYear, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}

// GetCurrent retrieves the current fiscal year for a tenant
func (s *fiscalYearService) GetCurrent(ctx context.Context, tenantID uuid.UUID) (*domain.FiscalYear, error) {
	return s.repo.GetCurrentByTenantID(ctx, tenantID)
}

// SetAsCurrent sets a fiscal year as current
func (s *fiscalYearService) SetAsCurrent(ctx context.Context, tenantID, fiscalYearID uuid.UUID) error {
	return s.repo.SetAsCurrent(ctx, tenantID, fiscalYearID)
}

// Close closes a fiscal year
func (s *fiscalYearService) Close(ctx context.Context, fiscalYearID, closedBy uuid.UUID) error {
	// Get fiscal year
	fy, err := s.repo.GetByID(ctx, fiscalYearID)
	if err != nil {
		return fmt.Errorf("failed to get fiscal year: %w", err)
	}

	// Check if already closed
	if fy.IsClosed {
		return fmt.Errorf("fiscal year is already closed")
	}

	// Close it
	fy.Close(closedBy)

	// Update in database
	if err := s.repo.Update(ctx, fy); err != nil {
		return fmt.Errorf("failed to close fiscal year: %w", err)
	}

	return nil
}

// Reopen reopens a closed fiscal year
func (s *fiscalYearService) Reopen(ctx context.Context, fiscalYearID uuid.UUID) error {
	// Get fiscal year
	fy, err := s.repo.GetByID(ctx, fiscalYearID)
	if err != nil {
		return fmt.Errorf("failed to get fiscal year: %w", err)
	}

	// Check if closed
	if !fy.IsClosed {
		return fmt.Errorf("fiscal year is not closed")
	}

	// Reopen it
	fy.Reopen()

	// Update in database
	if err := s.repo.Update(ctx, fy); err != nil {
		return fmt.Errorf("failed to reopen fiscal year: %w", err)
	}

	return nil
}

// Delete deletes a fiscal year
func (s *fiscalYearService) Delete(ctx context.Context, fiscalYearID uuid.UUID) error {
	// Get fiscal year
	fy, err := s.repo.GetByID(ctx, fiscalYearID)
	if err != nil {
		return fmt.Errorf("failed to get fiscal year: %w", err)
	}

	// Check if current
	if fy.IsCurrent {
		return fmt.Errorf("cannot delete current fiscal year")
	}

	// Check if closed
	if fy.IsClosed {
		return fmt.Errorf("cannot delete closed fiscal year")
	}

	// Delete
	if err := s.repo.Delete(ctx, fiscalYearID); err != nil {
		return fmt.Errorf("failed to delete fiscal year: %w", err)
	}

	return nil
}

// GenerateInvoiceNumber generates the next invoice number
func (s *fiscalYearService) GenerateInvoiceNumber(ctx context.Context, fiscalYearID uuid.UUID) (string, error) {
	// Get fiscal year
	fy, err := s.repo.GetByID(ctx, fiscalYearID)
	if err != nil {
		return "", fmt.Errorf("failed to get fiscal year: %w", err)
	}

	// Check if closed
	if fy.IsClosed {
		return "", fmt.Errorf("cannot generate invoice number for closed fiscal year")
	}

	// Increment number
	nextNum, err := s.repo.IncrementInvoiceNumber(ctx, fiscalYearID)
	if err != nil {
		return "", fmt.Errorf("failed to increment invoice number: %w", err)
	}

	// Generate full number (e.g., "INV-8283-0001")
	return fmt.Sprintf("%s%04d", fy.InvoicePrefix, nextNum), nil
}

// GeneratePurchaseNumber generates the next purchase number
func (s *fiscalYearService) GeneratePurchaseNumber(ctx context.Context, fiscalYearID uuid.UUID) (string, error) {
	// Get fiscal year
	fy, err := s.repo.GetByID(ctx, fiscalYearID)
	if err != nil {
		return "", fmt.Errorf("failed to get fiscal year: %w", err)
	}

	// Check if closed
	if fy.IsClosed {
		return "", fmt.Errorf("cannot generate purchase number for closed fiscal year")
	}

	// Increment number
	nextNum, err := s.repo.IncrementPurchaseNumber(ctx, fiscalYearID)
	if err != nil {
		return "", fmt.Errorf("failed to increment purchase number: %w", err)
	}

	// Generate full number (e.g., "PUR-8283-0001")
	return fmt.Sprintf("%s%04d", fy.PurchasePrefix, nextNum), nil
}

// GenerateVoucherNumber generates the next voucher number
func (s *fiscalYearService) GenerateVoucherNumber(ctx context.Context, fiscalYearID uuid.UUID) (string, error) {
	// Get fiscal year
	fy, err := s.repo.GetByID(ctx, fiscalYearID)
	if err != nil {
		return "", fmt.Errorf("failed to get fiscal year: %w", err)
	}

	// Check if closed
	if fy.IsClosed {
		return "", fmt.Errorf("cannot generate voucher number for closed fiscal year")
	}

	// Increment number
	nextNum, err := s.repo.IncrementVoucherNumber(ctx, fiscalYearID)
	if err != nil {
		return "", fmt.Errorf("failed to increment voucher number: %w", err)
	}

	// Generate full number (e.g., "JV-8283-0001")
	return fmt.Sprintf("%s%04d", fy.VoucherPrefix, nextNum), nil
}
