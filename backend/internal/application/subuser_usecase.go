package application

import (
	"app/internal/domain/user"
	"fmt"
	"time"
)

var PROVIDER_SECONDARY = 3

type SubUserUseCase struct {
	userRepo user.Repository
}

func NewSubUserUseCase(userRepo user.Repository) *SubUserUseCase {
	return &SubUserUseCase{
		userRepo: userRepo,
	}
}

// CreateSubUser creates a subuser for a given main user's username
func (uc *SubUserUseCase) CreateSubUser(mainUsername, subUsername, subPassword string) (*user.User, error) {
	// 1. Get main user by username
	mainUser, err := uc.userRepo.GetByLogin(mainUsername)
	if err != nil {
		if uc.userRepo.IsNotFoundError(err) {
			return nil, fmt.Errorf("main user not found: %s", mainUsername)
		}
		return nil, fmt.Errorf("error retrieving main user: %w", err)
	}

	// 2. Create subuser
	subUser := &user.User{
		Login:       subUsername,
		OwnerID:     &mainUser.ID,
		Active:      true,
		IsLogged:    false,
		CompanyID:   mainUser.CompanyID,
		CompanyName: mainUser.CompanyName,

		ProviderID: uint(PROVIDER_SECONDARY),
		CreatedAt:  time.Now().Format(time.RFC3339),
		LastAccess: time.Now().Format(time.RFC3339),
	}

	// Set password for subuser
	if err := subUser.SetPassword(subPassword); err != nil {
		return nil, fmt.Errorf("error setting password for subuser: %w", err)
	}

	// 3. Save subuser to the repository
	if err := uc.userRepo.Create(subUser); err != nil {
		return nil, fmt.Errorf("error creating subuser: %w", err)
	}

	subUser, err = uc.userRepo.GetByLogin(subUsername)
	if err != nil {
		return nil, fmt.Errorf("error retrieving subuser: %w", err)
	}

	return subUser, nil
}
