package application

import "app/internal/domain/user"

type UserUseCase struct {
	repo user.Repository
}

func NewUserUseCase(r user.Repository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) GetUserByID(id uint) (*user.User, error) {
	return uc.repo.GetByID(id)
}
