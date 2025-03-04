package user

type Repository interface {
	Create(u *User) error
	Update(u *User) error
	GetByID(id uint) (*User, error)
	GetByLogin(login string) (*User, error)
	GetByOwnerID(ownerID uint) ([]*User, error)

	// Methods for error handling check if the error is a not found error
	IsNotFoundError(err error) bool
}
