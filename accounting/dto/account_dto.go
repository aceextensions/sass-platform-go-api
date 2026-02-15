package dto

import (
	"github.com/aceextension/accounting/domain"
	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	Code        string             `json:"code" validate:"required"`
	Name        string             `json:"name" validate:"required"`
	Type        domain.AccountType `json:"type" validate:"required,oneof=ASSET LIABILITY EQUITY REVENUE EXPENSE"`
	ParentID    *uuid.UUID         `json:"parentId"`
	Description *string            `json:"description"`
}

type UpdateAccountRequest struct {
	Code        string             `json:"code"`
	Name        string             `json:"name"`
	Type        domain.AccountType `json:"type"`
	ParentID    *uuid.UUID         `json:"parentId"`
	IsActive    bool               `json:"isActive"`
	Description *string            `json:"description"`
}
