package application

import (
	"app/internal/domain/user"
	"fmt"
)

type ProfileUseCase struct {
	repo user.Repository
}

func (uc *ProfileUseCase) GetRepo() user.Repository {
	return uc.repo
}

func NewProfileUseCase(r user.Repository) *ProfileUseCase {
	return &ProfileUseCase{repo: r}
}

func (uc *ProfileUseCase) UploadProfile(ownerUsername string, login string, profile *user.Profile) (*user.User, error) {

	// Get user owner
	userOwner, err := uc.repo.GetByLogin(ownerUsername)
	if err != nil {
		if uc.repo.IsNotFoundError(err) {
			return nil, fmt.Errorf("user owner not found")
		}
		return nil, err
	}

	// Get user to update
	user, err := uc.repo.GetByLogin(login)
	if err != nil {
		return nil, err
	}

	if *user.OwnerID != userOwner.ID {
		return nil, fmt.Errorf("forbidden")
	}

	// Upload profile
	err = uc.repo.UploadProfileTransaction(user.ID, profile)
	if err != nil {
		return nil, err
	}

	user, err = uc.repo.GetByID(user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
