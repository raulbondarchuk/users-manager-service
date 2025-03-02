package user

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint
	UUID         string
	Login        string
	Password     string
	CompanyID    uint
	CompanyName  string
	ProviderID   uint
	ProviderName string
	Active       bool
	IsLogged     bool
	CreatedAt    string
	LastAccess   string
	Refresh      string
	RefreshExp   string
	OwnerID      uint
}

// Setter for password (hashing)
func (u *User) SetPassword(plain string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// Check password
func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
