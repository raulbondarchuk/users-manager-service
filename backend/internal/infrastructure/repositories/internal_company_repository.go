package repositories

import (
	"app/internal/domain/internal_company"
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/db/models"
	"errors"

	"gorm.io/gorm"
)

type internalCompanyRepository struct {
	db *gorm.DB
}

// Ensure providerRepository implements the domain interface
var _ internal_company.Repository = (*internalCompanyRepository)(nil)

func NewInternalCompanyRepository() internal_company.Repository {
	return &internalCompanyRepository{db: db.GetProvider().GetDB()}
}

func (r *internalCompanyRepository) IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (r *internalCompanyRepository) GetAll() ([]*internal_company.InternalCompany, error) {
	var companyModels []models.InternalCompanyModel
	if err := r.db.Find(&companyModels).Error; err != nil {
		return nil, err
	}

	companies := make([]*internal_company.InternalCompany, len(companyModels))
	for i, companyModel := range companyModels {
		companies[i] = companyModel.ToDomain()
	}
	return companies, nil
}

func (r *internalCompanyRepository) GetByName(name string) (*internal_company.InternalCompany, error) {
	var companyModel models.InternalCompanyModel
	if err := r.db.Where("name = ?", name).First(&companyModel).Error; err != nil {
		return nil, err
	}
	return companyModel.ToDomain(), nil
}

func (r *internalCompanyRepository) Create(companyName string) (*internal_company.InternalCompany, error) {
	companyModel := models.InternalCompanyModel{Name: companyName}
	if err := r.db.Create(&companyModel).Error; err != nil {
		return nil, err
	}
	return companyModel.ToDomain(), nil
}

func (r *internalCompanyRepository) GetByNameWithTransaction(tx *gorm.DB, name string) (*internal_company.InternalCompany, error) {
	var companyModel models.InternalCompanyModel
	if err := tx.Where("name = ?", name).First(&companyModel).Error; err != nil {
		return nil, err
	}
	return companyModel.ToDomain(), nil
}

func (r *internalCompanyRepository) CreateWithTransaction(tx *gorm.DB, companyName string) (*internal_company.InternalCompany, error) {
	companyModel := models.InternalCompanyModel{Name: companyName}
	if err := tx.Create(&companyModel).Error; err != nil {
		return nil, err
	}
	return companyModel.ToDomain(), nil
}
