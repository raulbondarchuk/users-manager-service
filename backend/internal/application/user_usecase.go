package application

import (
	"app/internal/domain/user"
)

type UserUseCase struct {
	repo user.Repository
}

func (uc *UserUseCase) GetRepo() user.Repository {
	return uc.repo
}

func NewUserUseCase(r user.Repository) *UserUseCase {
	return &UserUseCase{repo: r}
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
