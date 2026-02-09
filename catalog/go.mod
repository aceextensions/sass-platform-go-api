module github.com/aceextension/catalog

go 1.24.0

require (
	github.com/aceextension/audit v0.0.0
	github.com/aceextension/core v0.0.0
	github.com/aceextension/fiscal v0.0.0
	github.com/go-playground/validator/v10 v10.24.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/labstack/echo/v4 v4.13.3
)

replace (
	github.com/aceextension/audit => ../audit
	github.com/aceextension/core => ../core
	github.com/aceextension/fiscal => ../fiscal
)
