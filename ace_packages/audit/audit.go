package audit

import (
	"github.com/aceextension/audit/repository"
	"github.com/aceextension/audit/service"
)

// Global audit service instance
var Service service.AuditService

// Init initializes the audit module
func Init() {
	repo := repository.NewPostgresAuditRepository()
	Service = service.NewAuditService(repo)
}
