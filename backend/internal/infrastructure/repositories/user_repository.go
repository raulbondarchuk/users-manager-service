package repositories

import (
	"app/internal/domain/user"
	"app/internal/infrastructure/db"
	"app/internal/infrastructure/db/models"
	"app/pkg/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// Ensure userRepository implements the domain interface
var _ user.Repository = (*userRepository)(nil)

// NewUserRepository — constructor, accepts *gorm.DB
func NewUserRepository() user.Repository {
	return &userRepository{db: db.GetProvider().GetDB()}
}

func (r *userRepository) GetByID(id uint) (*user.User, error) {
	var um models.UserModel

	// Important to Preload("Profile"), so GORM loads the profile
	if err := r.db.Preload("Profile").Preload("Roles").First(&um, id).Error; err != nil {
		return nil, err
	}

	// Convert to domain entity
	domainUser := um.ToDomain()
	return domainUser, nil
}

func (r *userRepository) GetByLogin(login string) (*user.User, error) {
	var um models.UserModel
	err := r.db.Preload("Profile").Preload("Roles").
		Where("login = ?", login).
		First(&um).Error
	if err != nil {
		return nil, err
	}
	return um.ToDomain(), nil
}

func (r *userRepository) Create(u *user.User) error {
	um, err := db.FromDomainGeneric[user.User, models.UserModel](*u)
	if err != nil {
		return err
	}
	return r.db.Create(&um).Error
}

func (r *userRepository) Update(u *user.User) error {
	// Map to UserModel
	um, err := db.FromDomainGeneric[user.User, models.UserModel](*u)
	if err != nil {
		return err
	}
	return r.db.Save(&um).Error
}

func (r *userRepository) IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// GetByOwnerID returns users by OwnerID
func (r *userRepository) GetByOwnerID(ownerID uint) ([]*user.User, error) {
	var userModels []models.UserModel
	if err := r.db.Where("ownerId = ?", ownerID).Find(&userModels).Error; err != nil {
		return nil, err
	}

	users := make([]*user.User, len(userModels))
	for i, um := range userModels {
		users[i] = um.ToDomain()
	}
	return users, nil
}

// BeginTransaction starts a new transaction
func (r *userRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// CreateWithTransaction creates a user within a transaction
func (r *userRepository) CreateWithTransaction(tx *gorm.DB, userToCreate *user.User) error {
	um, err := db.FromDomainGeneric[user.User, models.UserModel](*userToCreate)
	if err != nil {
		return err
	}
	return tx.Create(&um).Error
}

// GetByLoginWithTransaction gets a user by login within a transaction
func (r *userRepository) GetByLoginWithTransaction(tx *gorm.DB, login string) (*user.User, error) {
	var userModel models.UserModel
	if err := tx.Preload("Roles").Where("login = ?", login).First(&userModel).Error; err != nil {
		return nil, err
	}
	return userModel.ToDomain(), nil
}

// GetUserAndSubUsersByIDWithTransaction retrieves a user by ID and all subusers with the same ownerId
func (r *userRepository) GetUserAndSubUsersByOwnerUsernameWithTransaction(tx *gorm.DB, ownerUsername string) (*user.User, []*user.User, error) {
	// Retrieve the main user by ID
	var mainUserModel models.UserModel
	if err := tx.Preload("Roles").Where("login = ?", ownerUsername).First(&mainUserModel).Error; err != nil {
		return nil, nil, err
	}

	// Convert main user model to domain entity
	mainUser := mainUserModel.ToDomain()

	// Retrieve all subusers with ownerId equal to the main user's ID
	var subUserModels []models.UserModel
	if err := tx.Where("ownerId = ?", mainUser.ID).Find(&subUserModels).Error; err != nil {
		return mainUser, nil, err
	}

	// Convert subuser models to domain entities
	subUsers := make([]*user.User, len(subUserModels))
	for i, um := range subUserModels {
		subUsers[i] = um.ToDomain()
	}

	return mainUser, subUsers, nil
}

func (r *userRepository) UpdateLastAccess(userId uint) error {
	return r.db.Model(&models.UserModel{}).Where("id = ?", userId).Update("lastAccess", time.Now().Format("2006-01-02 15:04:05")).Error
}

// UploadProfile updates the user's profile fields and lastAccess in a transaction
func (r *userRepository) UploadProfileTransaction(userId uint, profile *user.Profile) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"name":    utils.StringOrNil(profile.Name),
			"surname": utils.StringOrNil(profile.Surname),
			"email":   utils.StringOrNil(profile.Email),
			"phone":   utils.StringOrNil(profile.Phone),
			"role":    utils.StringOrNil(profile.Role),
			// "photo":   utils.StringOrNil(profile.Photo), // Uncomment if photo is a string
		}

		// Update the profile fields
		if err := tx.Model(&models.ProfileModel{}).Where("userId = ?", userId).Updates(updates).Error; err != nil {
			return err
		}

		// Update the lastAccess field of the associated user
		if err := tx.Model(&models.UserModel{}).Where("id = ?", userId).Update("lastAccess", time.Now().Format("2006-01-02 15:04:05")).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateActiveStatus updates the active status of a user
func (r *userRepository) UpdateActiveStatus(userID uint, active bool) error {
	return r.db.Model(&models.UserModel{}).Where("id = ?", userID).Update("active", active).Error
}
