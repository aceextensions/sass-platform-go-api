package dto

import "github.com/google/uuid"

type CreateWarehouseRequest struct {
	Name        string  `json:"name" validate:"required"`
	Location    string  `json:"location" validate:"required"`
	Description *string `json:"description"`
}

type UpdateWarehouseRequest struct {
	Name        string  `json:"name"`
	Location    string  `json:"location"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive"`
}

type WarehouseResponse struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenantId"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description *string   `json:"description"`
	IsActive    bool      `json:"isActive"`
}
