package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/tenancy/domain"
	"github.com/google/uuid"
)

type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	Update(ctx context.Context, tenant *domain.Tenant) error
	UpdateSettings(ctx context.Context, id uuid.UUID, settings interface{}) error
}

type PortalService interface {
	GetPortalInfo(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error)
	UpdateProfile(ctx context.Context, tenantID uuid.UUID, tenant *domain.Tenant) error
	ToggleModule(ctx context.Context, tenantID uuid.UUID, moduleCode string, enabled bool) error
	ConfigureDatabase(ctx context.Context, tenantID uuid.UUID, dbURL string) error
}

type portalService struct {
	repo         TenantRepository
	migrationSvc *MigrationService
}

func NewPortalService(repo TenantRepository, migrationSvc *MigrationService) PortalService {
	return &portalService{
		repo:         repo,
		migrationSvc: migrationSvc,
	}
}

func (s *portalService) GetPortalInfo(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	return s.repo.GetByID(ctx, tenantID)
}

func (s *portalService) UpdateProfile(ctx context.Context, tenantID uuid.UUID, update *domain.Tenant) error {
	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}

	// Update only allowed fields
	tenant.Name = update.Name
	tenant.Email = update.Email
	tenant.Phone = update.Phone
	tenant.Metadata = update.Metadata

	return s.repo.Update(ctx, tenant)
}

func (s *portalService) ToggleModule(ctx context.Context, tenantID uuid.UUID, moduleCode string, enabled bool) error {
	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}

	// Unmarshal settings
	var settings map[string]interface{}
	if tenant.Settings != nil {
		bytes, _ := json.Marshal(tenant.Settings)
		json.Unmarshal(bytes, &settings)
	} else {
		settings = make(map[string]interface{})
	}

	// Update module status
	modules, ok := settings["modules"].(map[string]interface{})
	if !ok {
		modules = make(map[string]interface{})
	}
	modules[moduleCode] = enabled
	settings["modules"] = modules

	return s.repo.UpdateSettings(ctx, tenantID, settings)
}

func (s *portalService) ConfigureDatabase(ctx context.Context, tenantID uuid.UUID, dbURL string) error {
	tenant, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}

	// 1. If a new DB URL is provided, run migrations first
	currentURL := ""
	if tenant.DatabaseURL != nil {
		currentURL = *tenant.DatabaseURL
	}
	if dbURL != "" && dbURL != currentURL {
		if err := s.migrationSvc.RunTenantMigration(ctx, dbURL); err != nil {
			return fmt.Errorf("failed to migrate dedicated database: %w", err)
		}
	}

	// 2. Update tenant record
	tenant.DatabaseURL = &dbURL
	err = s.repo.Update(ctx, tenant)
	if err != nil {
		return err
	}

	// 3. Hot-register the new pool in the DB Router
	if dbURL != "" {
		return db.Router.RegisterTenantPool(tenantID, dbURL)
	}

	return nil
}
