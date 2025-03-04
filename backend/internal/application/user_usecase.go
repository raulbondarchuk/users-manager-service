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

func (uc *UserUseCase) GetUserByLogin(login string) (*user.User, error) {
	return uc.repo.GetByLogin(login)
}

func (uc *UserUseCase) GetUserAndSubUsersByOwnerUsername(ownerUsername string) (*user.User, []*user.User, error) {

	tx := uc.repo.BeginTransaction()
	defer tx.Rollback()

	mainUser, subUsers, err := uc.repo.GetUserAndSubUsersByOwnerUsernameWithTransaction(tx, ownerUsername)
	if err != nil {
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
