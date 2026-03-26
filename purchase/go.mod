module github.com/aceextension/purchase

go 1.25.0

replace (
	github.com/aceextension/accounting => ../accounting
	github.com/aceextension/common => ../common
	github.com/aceextension/core => ../core
	github.com/aceextension/fiscal => ../fiscal
	github.com/aceextension/inventory => ../inventory
)

require (
	github.com/aceextension/accounting v0.0.0-00010101000000-000000000000
	github.com/aceextension/core v0.0.0
	github.com/aceextension/inventory v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.8.0
	github.com/labstack/echo/v4 v4.15.1
)

require (
	github.com/aceextension/fiscal v0.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)
