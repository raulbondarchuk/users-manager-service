package internal_company

import "gorm.io/gorm"

type Repository interface {
	GetAll() ([]*InternalCompany, error)
	GetByName(name string) (*InternalCompany, error)
	Create(companyName string) (*InternalCompany, error)

	GetByNameWithTransaction(tx *gorm.DB, name string) (*InternalCompany, error)
	CreateWithTransaction(tx *gorm.DB, companyName string) (*InternalCompany, error)

	IsNotFoundError(err error) bool
}
