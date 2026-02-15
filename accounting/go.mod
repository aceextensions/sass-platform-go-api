module github.com/aceextension/accounting

go 1.24.0

require (
	github.com/aceextension/core v0.0.0
	github.com/aceextension/identity v0.0.0
	github.com/aceextension/fiscal v0.0.0
	github.com/google/uuid v1.6.0
)

replace (
	github.com/aceextension/core => ../core
	github.com/aceextension/identity => ../identity
	github.com/aceextension/fiscal => ../fiscal
)
