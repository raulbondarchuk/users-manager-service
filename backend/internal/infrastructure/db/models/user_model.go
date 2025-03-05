package models

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"app/internal/domain/user"
)

type UserModel struct {
	ID       uint    `gorm:"column:id;primaryKey"`
	OwnerID  *uint   `gorm:"column:ownerId;default:null"`
	UUID     string  `gorm:"column:uuid;type:char(36);not null;uniqueIndex"`
	Login    string  `gorm:"column:login;size:255;not null;uniqueIndex"`
	Password *string `gorm:"column:password;size:255;"`

	// Set default values in GORM tags
	CompanyID   uint   `gorm:"column:companyId;not null;default:1"`
	CompanyName string `gorm:"column:companyName;size:255;not null;default:'Liftel'"`

	// FK of Provider
	ProviderID   uint   `gorm:"column:providerId;not null;default:1"`
	ProviderName string `gorm:"column:providerName;size:255;not null;default:'Liftel'"`

	Refresh    *string `gorm:"column:refresh;size:255,default:null"` // Refresh token
	RefreshExp string  `gorm:"column:refreshExp;type:DATETIME"`      // Refresh token expiration date
	Active     bool    `gorm:"column:active;default:true"`
	IsLogged   bool    `gorm:"column:isLogged;default:false"`
	LastAccess string  `gorm:"column:lastAccess;type:DATETIME"`
	CreatedAt  string  `gorm:"column:createdAt;type:datetime"`

	// GORM will load the Provider automatically
	Provider ProviderModel `gorm:"foreignKey:ProviderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Profile  *ProfileModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	// Many-to-many relationship through ref_user_role
	Roles []RoleModel `gorm:"many2many:ref_user_role;foreignKey:ID;joinForeignKey:user_id;References:ID;joinReferences:role_id"`
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

	// Generate UUID if empty
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}
	u.LastAccess = time.Now().Format("2006-01-02 15:04:05")

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

// AfterCreate is called after creating UserModel
func (u *UserModel) AfterCreate(tx *gorm.DB) (err error) {
	// Determine if the profile is primary and set the role if needed
	isPrimary := u.OwnerID == nil
	var role *string
	if isPrimary {
		defaultRole := "Company"
		role = &defaultRole
	}

	// Create a profile
	prof := ProfileModel{
		UserID:    u.ID,
		IsPrimary: isPrimary,
		Role:      role,
	}

	if err := tx.Create(&prof).Error; err != nil {
		return err
	}
	return nil
}

// ---- Mapping: model <-> domain entity ----

// ToDomain converts UserModel to domain entity user.User
func (um *UserModel) ToDomain() *user.User {
	domainUser := user.User{
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

	// Parse LastAccess if it's not empty
	if um.RefreshExp != "" {
		parsedTime, err := time.Parse(time.RFC3339, um.RefreshExp)
		if err != nil {
			log.Printf("Error parsing RefreshExp time: %v", err)
		} else {
			domainUser.RefreshExp = parsedTime.Format("2006-01-02 15:04:05")
		}
	}

	// If UserModel has a profile (Preload("Profile") loaded it),
	// then convert it to domain.Profile
	if um.Profile != nil {
		domainUser.Profile = um.Profile.ToDomain()
	}

	// Load roles of user
	for _, role := range um.Roles {
		domainUser.Roles = append(domainUser.Roles, *role.ToDomain())
	}

	return &domainUser
}

// hashPasswordIfNeeded hashes the password if it's not already hashed
func (u *UserModel) hashPasswordIfNeeded() error {

	if u.Password == nil {
		return nil
	}

	// If password is not hashed (does not start with "$2" and etc.)
	if !strings.HasPrefix(*u.Password, "$2") {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		*u.Password = string(hashed)
	}
	return nil
}
