# How to Create a New Module

This guide provides a step-by-step workflow for creating a new module in the AceExtension Go Monorepo. Follow this pattern to ensure consistency, testability, and adherence to our architecture.

## 1. Module Structure Setup

Navigate to the project root and create your module directory.

```bash
mkdir -p mymodule/cmd mymodule/domain mymodule/dto mymodule/handler mymodule/migrations mymodule/repository mymodule/service
```

### Initialize `go.mod`

Create `mymodule/go.mod`. Use local replacements for internal dependencies.

```go
module github.com/aceextension/mymodule

go 1.24.0

require (
    github.com/aceextension/core v0.0.0
    github.com/google/uuid v1.6.0
    github.com/labstack/echo/v4 v4.15.0
)

replace (
    github.com/aceextension/core => ../core
)
```

**Action:** Add your module to the root `go.work` file.
```bash
go work use mymodule
```

---

## 2. Domain Layer (`domain/`)

Define your core entities here. These are pure Go structs with JSON tags.

**File:** `mymodule/domain/entity.go`
```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type MyEntity struct {
    ID        uuid.UUID `json:"id"`
    TenantID  uuid.UUID `json:"tenantId"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}

func NewMyEntity(tenantID uuid.UUID, name string) *MyEntity {
    return &MyEntity{
        ID:        uuid.New(),
        TenantID:  tenantID,
        Name:      name,
        CreatedAt: time.Now(),
    }
}
```

---

## 3. DTO Layer (`dto/`)

Define request/response structures for your API. strict validation tags are required.

**File:** `mymodule/dto/request.go`
```go
package dto

type CreateEntityRequest struct {
    Name string `json:"name" validate:"required,min=3,max=100"`
}
```

---

## 4. Repository Layer (`repository/`)

**Pattern:** Define the Interface first, then implement the PostgreSQL version.

**File:** `mymodule/repository/interfaces.go`
```go
package repository

import (
    "context"
    "github.com/aceextension/mymodule/domain"
    "github.com/google/uuid"
)

type Repository interface {
    Create(ctx context.Context, entity *domain.MyEntity) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.MyEntity, error)
}
```

**File:** `mymodule/repository/postgres_repository.go`
```go
package repository

import (
    "github.com/aceextension/core/db"
    // ... imports
)

type postgresRepository struct {
    pool db.QueryExecutor
}

func NewPostgresRepository(pool db.QueryExecutor) Repository {
    return &postgresRepository{pool: pool}
}
// ... Implement methods using r.pool.Exec or r.pool.QueryRow
```

---

## 5. Service Layer (`service/`)

Contains business logic. It accepts DTOs, validates them, interacts with the domain/repository, and returns domain objects.

**File:** `mymodule/service/service.go`
```go
package service

import (
    "github.com/aceextension/mymodule/domain"
    "github.com/aceextension/mymodule/dto"
    "github.com/aceextension/mymodule/repository"
)

type Service interface {
    Create(ctx context.Context, tenantID uuid.UUID, req dto.CreateEntityRequest) (*domain.MyEntity, error)
}

type serviceImpl struct {
    repo repository.Repository
}

func NewService(repo repository.Repository) Service {
    return &serviceImpl{repo: repo}
}

func (s *serviceImpl) Create(ctx, tenantID, req) (*domain.MyEntity, error) {
    // 1. Business Validation logic
    // 2. Create Domain Object
    entity := domain.NewMyEntity(tenantID, req.Name)
    // 3. Persist
    return entity, s.repo.Create(ctx, entity)
}
```

---

## 6. Handler Layer (`handler/`)

Handles HTTP requests (Echo context), parses bodies, calls Service, and returns JSON.

**File:** `mymodule/handler/handler.go`
```go
package handler

import (
    "github.com/aceextension/core/db" // for GetTenantID
    "github.com/labstack/echo/v4"
)

type Handler struct {
    service service.Service
}

// Swagger annotations here...
func (h *Handler) Create(c echo.Context) error {
    tenantID, _ := db.GetTenantID(c.Request().Context())
    var req dto.CreateEntityRequest
    if err := c.Bind(&req); err != nil { ... }
    if err := c.Validate(&req); err != nil { ... } // Auto-validation

    res, err := h.service.Create(c.Request().Context(), tenantID, req)
    if err != nil { ... }
    return c.JSON(http.StatusCreated, res)
}
```

---

## 7. Routes Registration

**File:** `mymodule/handler/routes.go`
```go
func RegisterRoutes(e *echo.Group, h *Handler) {
    g := e.Group("/mymodule")
    g.POST("", h.Create)
}
```

---

## 8. Main Integration (`api/cmd/api/main.go`)

1.  Update `api/go.mod` with `replace github.com/aceextension/mymodule => ../mymodule`.
2.  Run `go mod tidy` in `api/`.
3.  In `main.go`:
    ```go
    import (
        "github.com/aceextension/mymodule"
        myHandler "github.com/aceextension/mymodule/handler"
    )

    // ... inside main()
    mymodule.Init()
    h := myHandler.NewHandler(mymodule.Service)
    myHandler.RegisterRoutes(api.Group("/v1"), h)
    ```

---

## Checklist
- [ ] `go.mod` created and replaced?
- [ ] Structs have `json` tags?
- [ ] DTOs have `validate` tags?
- [ ] Repository uses `db.QueryExecutor`?
- [ ] SQL Migration created?
- [ ] Routes registered in `main.go`?
