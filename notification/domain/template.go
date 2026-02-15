package domain

import (
	"time"

	"github.com/google/uuid"
)

// ChannelType defines the notification channel
type ChannelType string

const (
	// SMS channel
	ChannelSMS ChannelType = "SMS"
	// Email channel
	ChannelEmail ChannelType = "EMAIL"
	// WhatsApp channel
	ChannelWhatsApp ChannelType = "WHATSAPP"
	// InApp channel
	ChannelInApp ChannelType = "IN_APP"
)

// Template represents a notification template
type Template struct {
	ID        uuid.UUID   `json:"id"`
	TenantID  uuid.UUID   `json:"tenantId"`
	Code      string      `json:"code"`
	Channel   ChannelType `json:"channel"`
	Subject   *string     `json:"subject,omitempty"`
	Body      string      `json:"body"`
	IsActive  bool        `json:"isActive"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// NewTemplate creates a new template
func NewTemplate(tenantID uuid.UUID, code string, channel ChannelType, body string) *Template {
	return &Template{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Code:      code,
		Channel:   channel,
		Body:      body,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
