package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeEquity    AccountType = "EQUITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeExpense   AccountType = "EXPENSE"
)

type Account struct {
	ID          uuid.UUID   `json:"id"`
	TenantID    uuid.UUID   `json:"tenantId"`
	Code        string      `json:"code"`     // "1001", "2000"
	Name        string      `json:"name"`     // "Cash on Hand", "Sales Revenue"
	Type        AccountType `json:"type"`     // ASSET, LIABILITY, etc.
	ParentID    *uuid.UUID  `json:"parentId"` // For hierarchical CoA
	IsActive    bool        `json:"isActive"`
	Description *string     `json:"description"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

func NewAccount(tenantID uuid.UUID, code, name string, accType AccountType) *Account {
	return &Account{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Code:      code,
		Name:      name,
		Type:      accType,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (a *Account) Validate() error {
	if a.Code == "" {
		return errors.New("account code is required")
	}
	if a.Name == "" {
		return errors.New("account name is required")
	}
	if a.Type == "" {
		return errors.New("account type is required")
	}
	return nil
}
