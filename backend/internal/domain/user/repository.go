package user

type Repository interface {
	GetByID(id uint) (*User, error)
}
