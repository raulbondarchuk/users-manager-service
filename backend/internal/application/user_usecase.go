package application

import (
	"app/internal/application/ports"
	"app/internal/domain/internal_company"
	"app/internal/domain/user"
	"app/pkg/errorsLib"
	"fmt"
	"time"
)

type UserUseCase struct {
	repo                user.Repository
	internalCompanyRepo internal_company.Repository
	verificacionesSvc   ports.VerificacionesService
}

func (uc *UserUseCase) GetRepo() user.Repository {
	return uc.repo
}

func NewUserUseCase(r user.Repository, internalCompanyRepo internal_company.Repository, verificacionesSvc ports.VerificacionesService) *UserUseCase {
	return &UserUseCase{repo: r, internalCompanyRepo: internalCompanyRepo, verificacionesSvc: verificacionesSvc}
}

func (uc *UserUseCase) GetUserByID(id uint) (*user.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) CheckIfUserIsCompany(login string) (bool, error) {
	user, err := uc.repo.GetByLogin(login)
	if err != nil {
		return false, err
	}

	if user.OwnerID == nil || *user.OwnerID == 0 {
		return true, nil
	}

	return false, nil
}

func (uc *UserUseCase) CheckIfUserIsLogged(login string) (bool, error) {
	user, err := uc.repo.GetByLogin(login)
	if err != nil {
		return false, err
	}
	return user.IsLogged, nil
}

func (uc *UserUseCase) GetUserByLogin(login string) (*user.User, error) {
	return uc.repo.GetByLogin(login)
}

func (uc *UserUseCase) GetUserAndSubUsersByOwnerUsername(ownerUsername string) (*user.User, []*user.User, error) {

	tx := uc.repo.BeginTransaction()
	defer tx.Rollback()

	mainUser, subUsers, err := uc.repo.GetUserAndSubUsersByOwnerUsernameWithTransaction(tx, ownerUsername)
	if err != nil {
		if uc.repo.IsNotFoundError(err) {
			return nil, nil, errorsLib.ErrNotFound
		}
		return nil, nil, err
	}

	tx.Commit()

	return mainUser, subUsers, nil
}

func (uc *UserUseCase) ActivateDeactivateUser(username string, active bool) error {

	user, err := uc.repo.GetByLogin(username)
	if err != nil {
		return err
	}

	switch active {
	case true:
		err = uc.activateUser(user.ID)
	case false:
		err = uc.deactivateUser(user.ID)
	}

	return err
}

func (uc *UserUseCase) activateUser(userID uint) error {
	return uc.repo.UpdateActiveStatus(userID, true)
}

func (uc *UserUseCase) deactivateUser(userID uint) error {
	return uc.repo.UpdateActiveStatus(userID, false)
}

func (uc *UserUseCase) RegisterCompanyUser(username, password, companyName string) (*user.User, error) {
	tx := uc.repo.BeginTransaction()
	defer tx.Rollback()

	// Check if company already exists
	_, err := uc.internalCompanyRepo.GetByNameWithTransaction(tx, companyName)
	if err != nil {
		if !uc.internalCompanyRepo.IsNotFoundError(err) {
			return nil, fmt.Errorf("error checking if company exists: %w", err)
		}
	} else {
		return nil, fmt.Errorf("company already exists")
	}

	// Create company
	company, err := uc.internalCompanyRepo.CreateWithTransaction(tx, companyName)
	if err != nil {
		return nil, fmt.Errorf("error creating company: %w", err)
	}

	// Check if user already exists
	existingUser, err := uc.repo.GetByLoginWithTransaction(tx, username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user already exists")
	}

	// Check if company already exists in verificaciones
	exists, err := uc.verificacionesSvc.CheckIfUserExists(username)
	if err != nil {
		return nil, fmt.Errorf("error checking if company exists: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("company already exists")
	}

	// Create new user
	newUser := &user.User{
		Login:        username,
		ProviderID:   1, // Assuming 1 is the ID for liftel
		ProviderName: "Liftel",
		OwnerID:      nil,
		CompanyID:    company.ID,
		CompanyName:  company.Name,
		Active:       true,
		IsLogged:     false,
	}
	newUser.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	// Set password
	if err := newUser.SetPassword(password); err != nil {
		return nil, fmt.Errorf("error setting password: %w", err)
	}

	// Create user in the repository
	if err := uc.repo.CreateWithTransaction(tx, newUser); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// Retrieve the created user
	createdUser, err := uc.repo.GetByLogin(username)
	if err != nil {
		return nil, fmt.Errorf("error retrieving created user: %w", err)
	}

	return createdUser, nil
}
