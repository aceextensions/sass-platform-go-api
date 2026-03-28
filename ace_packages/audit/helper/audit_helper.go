package helper

import (
	"context"

	"github.com/aceextension/audit"
	"github.com/aceextension/audit/domain"
	"github.com/google/uuid"
)

// AuditHelper provides convenient methods for logging audit events
type AuditHelper struct{}

// LogUserAction logs a user-related action
func (h *AuditHelper) LogUserAction(ctx context.Context, action string, userID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := userID.String()
	return audit.Service.Log(ctx, action, "User", &entityID, details, auditCtx)
}

// LogTenantAction logs a tenant-related action
func (h *AuditHelper) LogTenantAction(ctx context.Context, action string, tenantID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := tenantID.String()
	return audit.Service.Log(ctx, action, "Tenant", &entityID, details, auditCtx)
}

// LogSaleAction logs a sale-related action
func (h *AuditHelper) LogSaleAction(ctx context.Context, action string, saleID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := saleID.String()
	return audit.Service.Log(ctx, action, "Sale", &entityID, details, auditCtx)
}

// LogPurchaseAction logs a purchase-related action
func (h *AuditHelper) LogPurchaseAction(ctx context.Context, action string, purchaseID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := purchaseID.String()
	return audit.Service.Log(ctx, action, "Purchase", &entityID, details, auditCtx)
}

// LogJournalAction logs a journal entry-related action
func (h *AuditHelper) LogJournalAction(ctx context.Context, action string, journalID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := journalID.String()
	return audit.Service.Log(ctx, action, "JournalEntry", &entityID, details, auditCtx)
}

// LogAccountAction logs an account-related action
func (h *AuditHelper) LogAccountAction(ctx context.Context, action string, accountID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := accountID.String()
	return audit.Service.Log(ctx, action, "Account", &entityID, details, auditCtx)
}

// LogSupplierAction logs a supplier-related action
func (h *AuditHelper) LogSupplierAction(ctx context.Context, action string, supplierID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := supplierID.String()
	return audit.Service.Log(ctx, action, "Supplier", &entityID, details, auditCtx)
}

// LogCustomerAction logs a customer-related action
func (h *AuditHelper) LogCustomerAction(ctx context.Context, action string, customerID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := customerID.String()
	return audit.Service.Log(ctx, action, "Customer", &entityID, details, auditCtx)
}

// LogProductAction logs a product-related action
func (h *AuditHelper) LogProductAction(ctx context.Context, action string, productID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := productID.String()
	return audit.Service.Log(ctx, action, "Product", &entityID, details, auditCtx)
}

// LogPaymentAction logs a payment-related action
func (h *AuditHelper) LogPaymentAction(ctx context.Context, action string, paymentID uuid.UUID, details any, auditCtx *domain.AuditContext) error {
	entityID := paymentID.String()
	return audit.Service.Log(ctx, action, "Payment", &entityID, details, auditCtx)
}

// NewAuditHelper creates a new audit helper
func NewAuditHelper() *AuditHelper {
	return &AuditHelper{}
}
