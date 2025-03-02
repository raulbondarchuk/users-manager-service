package repositories

import (
	"app/internal/domain/provider"
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/db/models"

	"gorm.io/gorm"
)

type providerRepository struct {
	db *gorm.DB
}

// Ensure providerRepository implements the domain interface
var _ provider.Repository = (*providerRepository)(nil)

// NewProviderRepository — constructor, accepts *gorm.DB
func NewProviderRepository() provider.Repository {
	return &providerRepository{db: db.GetProvider().GetDB()}
}

// GetAll returns all providers (map from ProviderModel to Provider)
func (r *providerRepository) GetAll() ([]provider.Provider, error) {
	var providerModels []models.ProviderModel
	if err := r.db.Find(&providerModels).Error; err != nil {
		return nil, err
	}

	// Map to domain entities
	providers := make([]provider.Provider, len(providerModels))
	for i, pm := range providerModels {
		p := pm.ToDomain() // *provider.Provider
		providers[i] = *p
	}
	return providers, nil
}

// GetByID returns a provider by ID (or an error if not found)
func (r *providerRepository) GetByID(id uint) (*provider.Provider, error) {
	var pm models.ProviderModel
	if err := r.db.First(&pm, id).Error; err != nil {
		return nil, err
	}
	return pm.ToDomain(), nil
}
