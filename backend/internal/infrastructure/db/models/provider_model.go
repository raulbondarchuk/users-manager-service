package models

import (
	"app/internal/domain/provider"

	"gorm.io/gorm"
)

// ProviderModel — GORM-model for the providers table
type ProviderModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null;uniqueIndex"`
	Desc string `gorm:"size:255"`
}

func (ProviderModel) TableName() string { return "providers" }

// ToDomain converts ProviderModel to domain entity provider.Provider
func (pm *ProviderModel) ToDomain() *provider.Provider {
	return &provider.Provider{
		ID:   pm.ID,
		Name: pm.Name,
		Desc: pm.Desc,
	}
}

// ----------------- Loading provider name -------------------

// LoadProviderName loads provider name by ID from the database
func LoadProviderName(tx *gorm.DB, providerID uint) (string, error) {
	if providerID == 0 {
		return "", nil
	}

	var p ProviderModel
	if err := tx.First(&p, providerID).Error; err != nil {
		return "", err
	}

	return p.Name, nil
}
