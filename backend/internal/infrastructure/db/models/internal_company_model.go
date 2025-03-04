package models

import "app/internal/domain/internal_company"

type InternalCompanyModel struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null;unique;size:255"` // UQ
}

func (InternalCompanyModel) TableName() string {
	return "internal_companies"
}

// ToDomain converts InternalCompanyModel to domain entity company.Company
func (ic *InternalCompanyModel) ToDomain() *internal_company.InternalCompany {
	return &internal_company.InternalCompany{
		ID:   ic.ID,
		Name: ic.Name,
	}
}
