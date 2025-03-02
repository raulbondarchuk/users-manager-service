package models

import (
	"app/internal/domain/user"
)

type ProfileModel struct {
	ID        uint    `gorm:"primaryKey;column:id"`
	UserID    uint    `gorm:"uniqueIndex;not null;column:userId"` // 1:1 (each profile is associated with one user)
	IsPrimary bool    `gorm:"not null;default:false;column:isPrimary"`
	Name      *string `gorm:"size:255,default:null;column:name"`
	Surname   *string `gorm:"size:255,default:null;column:surname"`
	Email     *string `gorm:"size:255,default:null;column:email"`
	Phone     *string `gorm:"size:255,default:null;column:phone"`
	Role      *string `gorm:"size:255,default:null;column:role"`
	Photo     *string `gorm:"size:255,default:null;column:photo"`

	// GORM 1:1 connection
	// constraint:OnDelete:CASCADE => if the user is deleted, the profile will also be deleted
	UserModel UserModel `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ProfileModel) TableName() string { return "profiles" }

// ToDomain converts the model to a domain entity
func (pm *ProfileModel) ToDomain() *user.Profile {
	return &user.Profile{
		ID:        pm.ID,
		UserID:    pm.UserID,
		IsPrimary: pm.IsPrimary,
		Name:      pm.Name,
		Surname:   pm.Surname,
		Email:     pm.Email,
		Phone:     pm.Phone,
		Role:      pm.Role,
		Photo:     pm.Photo,
	}
}
