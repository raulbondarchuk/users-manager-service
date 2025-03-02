package repositories

import (
	"app/internal/domain/user"
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/db/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// Ensure providerRepository implements the domain interface
var _ user.Repository = (*userRepository)(nil)

// NewProviderRepository — constructor, accepts *gorm.DB
func NewUserRepository() user.Repository {
	return &userRepository{db: db.GetProvider().GetDB()}
}

func (r *userRepository) GetByID(id uint) (*user.User, error) {
	var um models.UserModel

	// Important to Preload("Profile"), so GORM loads the profile
	if err := r.db.Preload("Profile").First(&um, id).Error; err != nil {
		return nil, err
	}

	// Convert to domain entity
	domainUser := um.ToDomain()
	return domainUser, nil
}
