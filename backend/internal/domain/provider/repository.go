package provider

type Repository interface {
	GetAll() ([]Provider, error)
	GetByID(id uint) (*Provider, error)
}
