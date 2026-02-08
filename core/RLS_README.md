# Row-Level Security (RLS) Middleware

Multi-tenant data isolation using PostgreSQL Row-Level Security.

## Features

- ✅ Automatic tenant data isolation
- ✅ Super admin bypass capability
- ✅ Transaction-scoped tenant context
- ✅ Echo middleware for HTTP requests
- ✅ Multiple tenant extraction methods (JWT, subdomain, header)

## Architecture

### Components

1. **RLS Manager** (`db/rls_manager.go`)
   - Manages RLS policies
   - Sets/clears tenant context
   - Provides transaction helpers

2. **Tenant Context** (`db/tenant_context.go`)
   - Go context utilities
   - Store/retrieve tenant/user IDs
   - Super admin flag management

3. **Tenant Middleware** (`middleware/tenant_middleware.go`)
   - Extract tenant from JWT/subdomain/header
   - Inject into request context
   - Enforce tenant requirements

## Usage

### 1. Initialize RLS

```go
import "github.com/aceextension/core/db"

func main() {
    db.Init(mainConnStr, auditConnStr)
    db.InitRLS()  // Initialize RLS manager
}
```

### 2. Use in HTTP Handlers

```go
import (
    "github.com/aceextension/core/middleware"
    "github.com/labstack/echo/v4"
)

func setupRoutes(e *echo.Echo) {
    // Apply tenant middleware globally
    e.Use(middleware.TenantMiddleware)
    
    // Routes that require tenant
    api := e.Group("/api")
    api.Use(middleware.RequireTenant)
    api.GET("/users", listUsers)
    
    // Super admin routes
    admin := e.Group("/admin")
    admin.Use(middleware.SuperAdminMiddleware)
    admin.GET("/tenants", listAllTenants)
}
```

### 3. Execute Queries with Tenant Context

```go
// Method 1: Using RLS helper
err := db.RLS.ExecuteWithTenantContext(ctx, tenantID, func(tx pgx.Tx) error {
    // All queries in this transaction are scoped to tenantID
    rows, err := tx.Query(ctx, "SELECT * FROM users")
    // ...
    return nil
})

// Method 2: Manual context setting
err := db.BeginFunc(ctx, func(tx pgx.Tx) error {
    // Set tenant context
    if err := db.RLS.SetTenantContext(ctx, tx, tenantID); err != nil {
        return err
    }
    
    // Query with tenant isolation
    rows, err := tx.Query(ctx, "SELECT * FROM users")
    // ...
    return nil
})

// Method 3: Super admin (bypass RLS)
err := db.RLS.ExecuteWithSuperAdminContext(ctx, func(tx pgx.Tx) error {
    // Can see ALL tenants' data
    rows, err := tx.Query(ctx, "SELECT * FROM users")
    // ...
    return nil
})
```

### 4. Use Context Utilities

```go
// Add to context
ctx = db.WithTenantID(ctx, tenantID)
ctx = db.WithUserID(ctx, userID)
ctx = db.WithSuperAdmin(ctx, true)

// Retrieve from context
tenantID, ok := db.GetTenantID(ctx)
userID, ok := db.GetUserID(ctx)
isSuperAdmin := db.IsSuperAdmin(ctx)
```

### 5. Extract Tenant in Handlers

```go
func listUsers(c echo.Context) error {
    // Get tenant ID from context
    tenantID, err := middleware.GetTenantIDFromContext(c)
    if err != nil {
        return err
    }
    
    // Use tenant ID
    users, err := userService.GetByTenantID(c.Request().Context(), tenantID)
    // ...
}
```

## Tenant Extraction Methods

The middleware tries to extract tenant ID in this order:

1. **JWT Claims** (recommended)
   - Extracts from `tenant_id` claim in JWT
   - Set by authentication middleware

2. **Subdomain**
   - Extracts from subdomain (e.g., `tenant1.example.com`)
   - Requires subdomain-to-tenant mapping

3. **Custom Header**
   - Extracts from `X-Tenant-ID` header
   - Useful for API clients

## Database Setup

### Enable RLS on Tables

```sql
-- Enable RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Create tenant isolation policy
CREATE POLICY tenant_isolation ON users
    USING (
        tenant_id = current_setting('app.current_tenant_id', true)::uuid
        OR current_setting('app.is_super_admin', true)::boolean = true
    );
```

### Helper Functions

```sql
-- Get current tenant ID
SELECT get_current_tenant_id();

-- Check if super admin
SELECT is_super_admin();
```

## Security Considerations

1. **Always use transactions**: RLS context is transaction-scoped
2. **Validate tenant ownership**: Verify user belongs to tenant before setting context
3. **Audit super admin access**: Log all super admin operations
4. **Test RLS policies**: Ensure policies work as expected
5. **Handle errors gracefully**: Don't leak tenant information in errors

## Best Practices

1. **Use `ExecuteWithTenantContext`**: Simplest and safest method
2. **Apply middleware globally**: Ensure all routes have tenant context
3. **Require tenant explicitly**: Use `RequireTenant` middleware for protected routes
4. **Log context changes**: Audit when tenant context is set/cleared
5. **Test isolation**: Verify tenants can't access each other's data

## Example: Complete Flow

```go
// 1. User logs in, JWT includes tenant_id
// 2. Tenant middleware extracts tenant_id from JWT
// 3. Adds to request context
// 4. Handler retrieves tenant_id
// 5. Executes query with tenant context
// 6. RLS ensures only tenant's data is returned

func createSale(c echo.Context) error {
    tenantID, _ := middleware.GetTenantIDFromContext(c)
    
    err := db.RLS.ExecuteWithTenantContext(c.Request().Context(), tenantID, func(tx pgx.Tx) error {
        // Insert sale - automatically scoped to tenant
        _, err := tx.Exec(ctx, "INSERT INTO sales (tenant_id, ...) VALUES ($1, ...)", tenantID)
        return err
    })
    
    return c.JSON(200, map[string]string{"status": "success"})
}
```

## Troubleshooting

**Problem**: Queries return no results  
**Solution**: Ensure tenant context is set before querying

**Problem**: RLS policies not working  
**Solution**: Check if RLS is enabled: `SELECT tablename, rowsecurity FROM pg_tables WHERE tablename = 'users';`

**Problem**: Super admin can't see all data  
**Solution**: Ensure `app.is_super_admin` is set to `true` in transaction

## Testing

Run the example:
```bash
go run core/examples/rls_example.go
```

Expected output:
- Tenant 1 sees only their data
- Tenant 2 sees only their data
- Super admin sees all data
