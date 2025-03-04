package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(u *User) error
	Update(u *User) error
	GetByID(id uint) (*User, error)
	GetByLogin(login string) (*User, error)
	GetByOwnerID(ownerID uint) ([]*User, error)

	BeginTransaction() *gorm.DB
	CreateWithTransaction(tx *gorm.DB, user *User) error
	GetByLoginWithTransaction(tx *gorm.DB, login string) (*User, error)
	GetUserAndSubUsersByOwnerUsernameWithTransaction(tx *gorm.DB, ownerUsername string) (*User, []*User, error)

	// Methods for error handling check if the error is a not found error
	IsNotFoundError(err error) bool
}
