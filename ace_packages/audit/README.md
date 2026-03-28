# Audit Module

Centralized audit logging system for tracking all operations across the platform.

## Features

- ✅ Immutable audit trail
- ✅ Async logging (non-blocking)
- ✅ Sync logging (when needed)
- ✅ Tenant-scoped audit logs
- ✅ JSONB details for flexible metadata
- ✅ Search and filter capabilities
- ✅ Separate audit database

## Usage

### Initialize Audit Module

```go
import "github.com/aceextension/audit"

func main() {
    // Initialize audit module
    audit.Init()
}
```

### Log Audit Events

```go
import (
    "github.com/aceextension/audit"
    "github.com/aceextension/audit/domain"
)

// Create audit context
auditCtx := &domain.AuditContext{
    TenantID:  &tenantID,
    UserID:    &userID,
    IPAddress: &ipAddress,
    UserAgent: &userAgent,
}

// Log an action (async - recommended)
entityID := user.ID.String()
audit.Service.Log(ctx, "CREATE_USER", "User", &entityID, map[string]interface{}{
    "userName": user.Name,
    "role":     user.Role,
}, auditCtx)

// Log an action (sync - when you need to ensure it's written)
audit.Service.LogSync(ctx, "DELETE_USER", "User", &entityID, map[string]interface{}{
    "userName": user.Name,
}, auditCtx)
```

### Using Audit Helper

```go
import "github.com/aceextension/audit/helper"

helper := helper.NewAuditHelper()

// Log user action
helper.LogUserAction(ctx, "UPDATE_USER", userID, map[string]interface{}{
    "oldRole": "user",
    "newRole": "admin",
}, auditCtx)

// Log sale action
helper.LogSaleAction(ctx, "CREATE_SALE", saleID, map[string]interface{}{
    "invoiceNumber": "INV-001",
    "totalAmount":   1000.00,
}, auditCtx)
```

### Query Audit Logs

```go
// Get audit logs for a tenant
logs, err := audit.Service.GetByTenantID(ctx, tenantID, 50, 0)

// Get audit logs for an entity
logs, err := audit.Service.GetByEntity(ctx, "Sale", saleID.String(), 50, 0)

// Search with filters
filters := &repository.AuditSearchFilters{
    TenantID:  &tenantID,
    Action:    strPtr("CREATE_SALE"),
    StartDate: strPtr("2026-01-01"),
    EndDate:   strPtr("2026-12-31"),
    Limit:     100,
}
logs, err := audit.Service.Search(ctx, filters)
```

## Common Audit Actions

### User Management
- `CREATE_USER`, `UPDATE_USER`, `DELETE_USER`
- `ACTIVATE_USER`, `DEACTIVATE_USER`
- `CHANGE_USER_ROLE`

### Authentication
- `LOGIN`, `LOGIN_FAILED`
- `LOGOUT`
- `CHANGE_PASSWORD`, `RESET_PASSWORD`

### Sales
- `CREATE_SALE`, `UPDATE_SALE`, `DELETE_SALE`
- `POST_SALE`, `CANCEL_SALE`

### Purchases
- `CREATE_PURCHASE`, `UPDATE_PURCHASE`, `DELETE_PURCHASE`
- `POST_PURCHASE`, `CANCEL_PURCHASE`

### Accounting
- `CREATE_JOURNAL_ENTRY`, `POST_JOURNAL_ENTRY`
- `CREATE_ACCOUNT`, `UPDATE_ACCOUNT`
- `CLOSE_FISCAL_YEAR`

### Admin Operations
- `SUSPEND_TENANT`, `ACTIVATE_TENANT`
- `UPDATE_SUBSCRIPTION`
- `VERIFY_KYB`, `REJECT_KYB`

## Database Schema

The audit logs are stored in a separate database (`aceextension_audit`) with the following schema:

```sql
CREATE TABLE audit_logs (
    id UUID NOT NULL,
    tenant_id UUID,
    user_id UUID,
    action VARCHAR(100) NOT NULL,
    entity VARCHAR(100) NOT NULL,
    entity_id TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    details JSONB,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id, created_at)
);
```

## Best Practices

1. **Use Async Logging**: Use `Log()` for most cases to avoid blocking
2. **Use Sync Logging**: Use `LogSync()` only for critical operations
3. **Include Context**: Always provide tenant, user, IP, and user agent
4. **Meaningful Details**: Include relevant metadata in the details field
5. **Consistent Naming**: Use consistent action names across the platform
