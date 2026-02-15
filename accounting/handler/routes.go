package handler

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, accountHandler *AccountHandler, journalHandler *JournalHandler, reportHandler *ReportHandler) {
	accountingGroup := e.Group("/accounting")

	// Accounts
	accountingGroup.POST("/accounts", accountHandler.CreateAccount)
	accountingGroup.GET("/accounts", accountHandler.ListAccounts)
	accountingGroup.GET("/accounts/:id", accountHandler.GetAccount)
	accountingGroup.PUT("/accounts/:id", accountHandler.UpdateAccount)

	// Journal Entries
	accountingGroup.POST("/journals", journalHandler.CreateJournalEntry)
	accountingGroup.GET("/journals", journalHandler.ListJournalEntries)
	accountingGroup.GET("/journals/:id", journalHandler.GetJournalEntry)
	accountingGroup.POST("/journals/:id/post", journalHandler.PostJournalEntry)

	// Reports
	accountingGroup.GET("/reports/general-ledger", reportHandler.GetGeneralLedger)
}
