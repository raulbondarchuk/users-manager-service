package user

import (
	"app/internal/domain/role"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint    `json:"id"`
	UUID         string  `json:"uuid"`
	Login        string  `json:"login"`
	Password     *string `json:"-"`
	CompanyID    uint    `json:"companyId"`
	CompanyName  string  `json:"companyName"`
	ProviderID   uint    `json:"providerId"`
	ProviderName string  `json:"providerName"`
	Active       bool    `json:"active"`
	IsLogged     bool    `json:"isLogged"`
	CreatedAt    string  `json:"createdAt"`
	LastAccess   string  `json:"lastAccess"`

	Refresh    *string `json:"refresh"`
	RefreshExp string  `json:"refreshExp"`
	OwnerID    *uint   `json:"ownerId"`

	Profile *Profile    `json:"profile"`
	Roles   []role.Role `json:"roles"`

	AccessToken   string `json:"-"`
	OwnerUsername string `json:"-"`
}

// Setter for password (hashing)
func (u *User) SetPassword(plain string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedStr := string(hashed)
	u.Password = &hashedStr
	return nil
}

// Check password
func (u *User) CheckPassword(plain string) bool {
	if u.Password == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(plain))
	return err == nil
}

type Profile struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"userId"`    // 1:1 connection with User
	IsPrimary bool    `json:"isPrimary"` // depends on whether the user has an OwnerID
	Name      *string `json:"name"`
	Surname   *string `json:"surname"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Role      *string `json:"role"`
	Photo     *string `json:"photo"` // link to photo (logo of profile)
}
