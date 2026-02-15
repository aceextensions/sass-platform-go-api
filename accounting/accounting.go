package accounting

import (
	"log"

	"github.com/aceextension/accounting/repository"
	"github.com/aceextension/accounting/service"
	"github.com/aceextension/core/db"
	"github.com/aceextension/fiscal"
)

var (
	Service service.AccountingService
)

func Init() {
	log.Println("Initializing Accounting Module...")
	repoAccount := repository.NewPostgresAccountRepository(db.MainPool)
	repoJournal := repository.NewPostgresJournalRepository(db.MainPool)

	Service = service.NewAccountingService(repoAccount, repoJournal, fiscal.Service)
	log.Println("Accounting Module Initialized")
}
