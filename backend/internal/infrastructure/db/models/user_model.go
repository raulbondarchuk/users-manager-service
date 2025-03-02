package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"app/internal/domain/user"
)

type UserModel struct {
	ID       uint   `gorm:"column:id;primaryKey"`
	OwnerID  uint   `gorm:"column:ownerId;default:null"`
	UUID     string `gorm:"column:uuid;type:char(36);not null;uniqueIndex"`
	Login    string `gorm:"column:login;size:255;not null;uniqueIndex"`
	Password string `gorm:"column:password;size:255;not null"`

	// Set default values in GORM tags
	CompanyID   uint   `gorm:"column:companyId;not null;default:1"`
	CompanyName string `gorm:"column:companyName;size:255;not null;default:'Liftel'"`

	// FK of Provider
	ProviderID   uint   `gorm:"column:providerId;not null;default:1"`
	ProviderName string `gorm:"column:providerName;size:255;not null;default:'Liftel'"`

	Active     bool    `gorm:"column:active;default:true"`
	CreatedAt  string  `gorm:"column:createdAt;type:datetime"`
	LastAccess string  `gorm:"column:lastAccess;type:DATETIME"`
	Refresh    *string `gorm:"column:refresh;size:255,default:null"` // Refresh token
	RefreshExp string  `gorm:"column:refreshExp;type:DATETIME"`      // Refresh token expiration date
	IsLogged   bool    `gorm:"column:isLogged;default:false"`

	// GORM will load the Provider automatically
	Provider ProviderModel `gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (UserModel) TableName() string { return "users" }

// ----------------- GORM callbacks ------------------

// BeforeCreate is called before creating a record (GORM callback).
func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate UUID if empty
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}

	u.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	u.LastAccess = time.Now().Format("2006-01-02 15:04:05")
	u.RefreshExp = time.Now().Format("2006-01-02 15:04:05")

	// Hash password if needed
	if err := u.hashPasswordIfNeeded(); err != nil {
		return err
	}

	// Update ProviderName from the provider table if ProviderID > 0
	providerName, err := LoadProviderName(tx, u.ProviderID)
	if err != nil {
		return err
	}
	u.ProviderName = providerName

	return nil
}

// BeforeSave is called before saving (create + update).
func (u *UserModel) BeforeSave(tx *gorm.DB) (err error) {
	return nil
}

// ---- Mapping: model <-> domain entity ----

// ToDomain converts UserModel to domain entity user.User
func (um *UserModel) ToDomain() *user.User {
	return &user.User{
		ID:           um.ID,
		UUID:         um.UUID,
		Login:        um.Login,
		Password:     um.Password,
		CompanyID:    um.CompanyID,
		CompanyName:  um.CompanyName,
		ProviderID:   um.ProviderID,
		ProviderName: um.ProviderName,
		Active:       um.Active,
		IsLogged:     um.IsLogged,
		CreatedAt:    um.CreatedAt,
		LastAccess:   um.LastAccess,
		Refresh:      um.Refresh,
		RefreshExp:   um.RefreshExp,
		OwnerID:      um.OwnerID,
	}
}

// hashPasswordIfNeeded hashes the password if it's not already hashed
func (u *UserModel) hashPasswordIfNeeded() error {
	// If password is not hashed (does not start with "$2" and etc.)
	if !strings.HasPrefix(u.Password, "$2") {
		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashed)
	}
	return nil
}
