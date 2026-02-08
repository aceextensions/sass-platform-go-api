package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/aceextension/audit"
	auditDomain "github.com/aceextension/audit/domain"
	crmDomain "github.com/aceextension/crm/domain"
	"github.com/aceextension/crm/repository"
	"github.com/aceextension/fiscal"
	"github.com/google/uuid"
)

// SupplierService defines the interface for supplier operations
type SupplierService interface {
	Create(ctx context.Context, supplier *crmDomain.Supplier) error
	GetByID(ctx context.Context, id uuid.UUID) (*crmDomain.Supplier, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*crmDomain.Supplier, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*crmDomain.Supplier, error)
	Update(ctx context.Context, supplier *crmDomain.Supplier) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*crmDomain.Supplier, error)
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// supplierService implements SupplierService
type supplierService struct {
	repo repository.SupplierRepository
}

// NewSupplierService creates a new supplier service
func NewSupplierService(repo repository.SupplierRepository) SupplierService {
	return &supplierService{
		repo: repo,
	}
}

// Create creates a new supplier
func (s *supplierService) Create(ctx context.Context, supplier *crmDomain.Supplier) error {
	// Generate supplier code if not provided
	if supplier.SupplierCode == "" {
		code, err := s.generateSupplierCode(ctx, supplier.TenantID)
		if err != nil {
			return fmt.Errorf("failed to generate supplier code: %w", err)
		}
		supplier.SupplierCode = code
	}

	// Create supplier
	if err := s.repo.Create(ctx, supplier); err != nil {
		return fmt.Errorf("failed to create supplier: %w", err)
	}

	// Audit log
	userID := uuid.Nil
	auditCtx := &auditDomain.AuditContext{
		UserID:   &userID, // TODO: Get from context
		TenantID: &supplier.TenantID,
	}

	entityIDStr := supplier.ID.String()
	audit.Service.Log(ctx, "CREATE_SUPPLIER", "Supplier", &entityIDStr, map[string]interface{}{
		"supplier_code": supplier.SupplierCode,
		"name":          supplier.Name,
		"email":         supplier.Email,
	}, auditCtx)

	return nil
}

// GetByID retrieves a supplier by ID
func (s *supplierService) GetByID(ctx context.Context, id uuid.UUID) (*crmDomain.Supplier, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCode retrieves a supplier by supplier code
func (s *supplierService) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*crmDomain.Supplier, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// GetByTenantID retrieves all suppliers for a tenant
func (s *supplierService) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*crmDomain.Supplier, error) {
	return s.repo.GetByTenantID(ctx, tenantID, limit, offset)
}

// Update updates a supplier
func (s *supplierService) Update(ctx context.Context, supplier *crmDomain.Supplier) error {
	// Get old supplier for audit
	oldSupplier, err := s.repo.GetByID(ctx, supplier.ID)
	if err != nil {
		return fmt.Errorf("failed to get old supplier: %w", err)
	}

	// Update supplier
	if err := s.repo.Update(ctx, supplier); err != nil {
		return fmt.Errorf("failed to update supplier: %w", err)
	}

	// Audit log
	userID := uuid.Nil
	auditCtx := &auditDomain.AuditContext{
		UserID:   &userID, // TODO: Get from context
		TenantID: &supplier.TenantID,
	}

	entityIDStr := supplier.ID.String()
	audit.Service.Log(ctx, "UPDATE_SUPPLIER", "Supplier", &entityIDStr, map[string]interface{}{
		"old_name":  oldSupplier.Name,
		"new_name":  supplier.Name,
		"old_email": oldSupplier.Email,
		"new_email": supplier.Email,
	}, auditCtx)

	return nil
}

// Delete deletes a supplier
func (s *supplierService) Delete(ctx context.Context, id uuid.UUID) error {
	// Get supplier for audit
	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get supplier: %w", err)
	}

	// Delete supplier
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	// Audit log
	userID := uuid.Nil
	auditCtx := &auditDomain.AuditContext{
		UserID:   &userID, // TODO: Get from context
		TenantID: &supplier.TenantID,
	}

	entityIDStr := id.String()
	audit.Service.Log(ctx, "DELETE_SUPPLIER", "Supplier", &entityIDStr, map[string]interface{}{
		"supplier_code": supplier.SupplierCode,
		"name":          supplier.Name,
	}, auditCtx)

	return nil
}

// Search searches suppliers
func (s *supplierService) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*crmDomain.Supplier, error) {
	return s.repo.Search(ctx, tenantID, query, limit, offset)
}

// Count returns total number of suppliers
func (s *supplierService) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, tenantID)
}

// generateSupplierCode generates a supplier code with fiscal year
func (s *supplierService) generateSupplierCode(ctx context.Context, tenantID uuid.UUID) (string, error) {
	// Get current fiscal year
	currentFY, err := fiscal.Service.GetCurrent(ctx, tenantID)
	if err != nil {
		// If no fiscal year, use simple numbering
		nextNum, err := s.repo.GetNextSupplierNumber(ctx, tenantID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("SUPP-%04d", nextNum), nil
	}

	// Generate with fiscal year (e.g., SUPP-8283-0001)
	nextNum, err := s.repo.GetNextSupplierNumber(ctx, tenantID)
	if err != nil {
		return "", err
	}

	yearCode := strings.ReplaceAll(currentFY.Name, "/", "")
	return fmt.Sprintf("SUPP-%s-%04d", yearCode, nextNum), nil
}
