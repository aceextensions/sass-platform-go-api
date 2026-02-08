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

// CustomerService defines the interface for customer operations
type CustomerService interface {
	Create(ctx context.Context, customer *crmDomain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*crmDomain.Customer, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*crmDomain.Customer, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*crmDomain.Customer, error)
	Update(ctx context.Context, customer *crmDomain.Customer) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*crmDomain.Customer, error)
	Count(ctx context.Context, tenantID uuid.UUID) (int64, error)
}

// customerService implements CustomerService
type customerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(repo repository.CustomerRepository) CustomerService {
	return &customerService{
		repo: repo,
	}
}

// Create creates a new customer
func (s *customerService) Create(ctx context.Context, customer *crmDomain.Customer) error {
	// Generate customer code if not provided
	if customer.CustomerCode == "" {
		code, err := s.generateCustomerCode(ctx, customer.TenantID)
		if err != nil {
			return fmt.Errorf("failed to generate customer code: %w", err)
		}
		customer.CustomerCode = code
	}

	// Create customer
	if err := s.repo.Create(ctx, customer); err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	// Audit log
	userID := uuid.Nil
	auditCtx := &auditDomain.AuditContext{
		UserID:   &userID, // TODO: Get from context
		TenantID: &customer.TenantID,
	}

	entityIDStr := customer.ID.String()
	audit.Service.Log(ctx, "CREATE_CUSTOMER", "Customer", &entityIDStr, map[string]interface{}{
		"customer_code": customer.CustomerCode,
		"name":          customer.Name,
		"email":         customer.Email,
	}, auditCtx)

	return nil
}

// GetByID retrieves a customer by ID
func (s *customerService) GetByID(ctx context.Context, id uuid.UUID) (*crmDomain.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByCode retrieves a customer by customer code
func (s *customerService) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*crmDomain.Customer, error) {
	return s.repo.GetByCode(ctx, tenantID, code)
}

// GetByTenantID retrieves all customers for a tenant
func (s *customerService) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*crmDomain.Customer, error) {
	return s.repo.GetByTenantID(ctx, tenantID, limit, offset)
}

// Update updates a customer
func (s *customerService) Update(ctx context.Context, customer *crmDomain.Customer) error {
	// Get old customer for audit
	oldCustomer, err := s.repo.GetByID(ctx, customer.ID)
	if err != nil {
		return fmt.Errorf("failed to get old customer: %w", err)
	}

	// Update customer
	if err := s.repo.Update(ctx, customer); err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	// Audit log
	userID := uuid.Nil
	auditCtx := &auditDomain.AuditContext{
		UserID:   &userID, // TODO: Get from context
		TenantID: &customer.TenantID,
	}

	entityIDStr := customer.ID.String()
	audit.Service.Log(ctx, "UPDATE_CUSTOMER", "Customer", &entityIDStr, map[string]interface{}{
		"old_name":  oldCustomer.Name,
		"new_name":  customer.Name,
		"old_email": oldCustomer.Email,
		"new_email": customer.Email,
	}, auditCtx)

	return nil
}

// Delete deletes a customer
func (s *customerService) Delete(ctx context.Context, id uuid.UUID) error {
	// Get customer for audit
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	// Delete customer
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	// Audit log
	userID := uuid.Nil
	auditCtx := &auditDomain.AuditContext{
		UserID:   &userID, // TODO: Get from context
		TenantID: &customer.TenantID,
	}

	entityIDStr := id.String()
	audit.Service.Log(ctx, "DELETE_CUSTOMER", "Customer", &entityIDStr, map[string]interface{}{
		"customer_code": customer.CustomerCode,
		"name":          customer.Name,
	}, auditCtx)

	return nil
}

// Search searches customers
func (s *customerService) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*crmDomain.Customer, error) {
	return s.repo.Search(ctx, tenantID, query, limit, offset)
}

// Count returns total number of customers
func (s *customerService) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	return s.repo.Count(ctx, tenantID)
}

// generateCustomerCode generates a customer code with fiscal year
func (s *customerService) generateCustomerCode(ctx context.Context, tenantID uuid.UUID) (string, error) {
	// Get current fiscal year
	currentFY, err := fiscal.Service.GetCurrent(ctx, tenantID)
	if err != nil {
		// If no fiscal year, use simple numbering
		nextNum, err := s.repo.GetNextCustomerNumber(ctx, tenantID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("CUST-%04d", nextNum), nil
	}

	// Generate with fiscal year (e.g., CUST-8283-0001)
	nextNum, err := s.repo.GetNextCustomerNumber(ctx, tenantID)
	if err != nil {
		return "", err
	}

	yearCode := strings.ReplaceAll(currentFY.Name, "/", "")
	return fmt.Sprintf("CUST-%s-%04d", yearCode, nextNum), nil
}
