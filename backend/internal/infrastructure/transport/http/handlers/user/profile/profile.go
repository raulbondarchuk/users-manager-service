package profile

import (
	"app/internal/application"
	"app/internal/domain/user"
	"app/internal/infrastructure/token/paseto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileUC *application.ProfileUseCase
}

func NewProfileHandler(profileUC *application.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{profileUC: profileUC}
}

type ProfileRequest struct {
	Name    *string `json:"name"`
	Surname *string `json:"surname"`
	Email   *string `json:"email"`
	Phone   *string `json:"phone"`
	Role    *string `json:"role"`
	Photo   *string `json:"photo"`
}

func (h *ProfileHandler) UpdateOwnProfile(c *gin.Context) {
	claims, err := paseto.Paseto().ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	h.updateProfile(c, claims.Username)
}

// UpdateUserProfile updates the profile of a specific user by user ID
func (h *ProfileHandler) UpdateUserProfile(c *gin.Context) {
	username := c.Query("username")
	h.updateProfile(c, username)
}

// UpdateUserProfile updates the profile of a specific user by user ID
func (h *ProfileHandler) updateProfile(c *gin.Context, username string) {

	claims, err := paseto.Paseto().ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	var ownerUsername string
	if claims.OwnerUsername == "" {
		ownerUsername = claims.Username
	} else {
		ownerUsername = claims.OwnerUsername
	}

	var profile ProfileRequest
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile request"})
		return
	}

	profileUseCase := user.Profile{
		Name:    profile.Name,
		Surname: profile.Surname,
		Email:   profile.Email,
		Phone:   profile.Phone,
		Role:    profile.Role,
		Photo:   profile.Photo,
	}

	user, err := h.profileUC.UploadProfile(ownerUsername, username, &profileUseCase)
	if err != nil {
		if h.profileUC.GetRepo().IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not allowed to update this profile"})
		} else if err.Error() == "user owner not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "user owner not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}
