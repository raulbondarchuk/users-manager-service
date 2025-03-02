package application

import (
	"app/internal/domain/provider"
)

type ProviderUseCase struct {
	repo provider.Repository
}

func NewProviderUseCase(r provider.Repository) *ProviderUseCase {
	return &ProviderUseCase{repo: r}
}

func (uc *ProviderUseCase) GetAllProviders() ([]provider.Provider, error) {
	return uc.repo.GetAll()
}

func (uc *ProviderUseCase) GetProviderByID(id uint) (*provider.Provider, error) {
	return uc.repo.GetByID(id)
}
