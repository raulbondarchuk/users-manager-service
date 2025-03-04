package application

import "app/internal/domain/user"

type ProfileUseCase struct {
	repo user.Repository
}

func (uc *ProfileUseCase) GetRepo() user.Repository {
	return uc.repo
}

func NewProfileUseCase(r user.Repository) *ProfileUseCase {
	return &ProfileUseCase{repo: r}
}

func (uc *ProfileUseCase) UploadProfile(login string, profile *user.Profile) (*user.User, error) {

	user, err := uc.repo.GetByLogin(login)
	if err != nil {
		return nil, err
	}

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
