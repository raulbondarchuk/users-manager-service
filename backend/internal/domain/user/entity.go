package user

import (
	"app/internal/domain/role"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint    `json:"id"`
	UUID         string  `json:"-"` // `json:"uuid"`
	Login        string  `json:"username"`
	Password     *string `json:"-"`
	CompanyID    uint    `json:"-"` //`json:"companyId"`
	CompanyName  string  `json:"-"` //`json:"companyName"`
	ProviderID   uint    `json:"providerId"`
	ProviderName string  `json:"providerName"`
	Active       bool    `json:"active"`
	IsLogged     bool    `json:"isLogged"`
	CreatedAt    string  `json:"createdAt"`
	LastAccess   string  `json:"lastAccess"`

	Refresh    *string `json:"-"` // `json:"refresh"`
	RefreshExp string  `json:"-"` // `json:"refreshExp"`
	OwnerID    *uint   `json:"-"` // `json:"ownerId"`

	Profile *Profile    `json:"profile"`
	Roles   []role.Role `json:"-"` // `json:"roles"`

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
	ID        uint    `json:"-"`         //`json:"id"`
	UserID    uint    `json:"-"`         //`json:"userId"`    // 1:1 connection with User
	IsPrimary bool    `json:"isPrimary"` // depends on whether the user has an OwnerID
	Name      *string `json:"name"`
	Surname   *string `json:"surname"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Role      *string `json:"role"`
	Photo     *string `json:"photo"` // link to photo (logo of profile)
}
