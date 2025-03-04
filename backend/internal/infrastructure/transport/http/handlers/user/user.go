package user

import (
	"app/internal/application"
	"app/internal/infrastructure/token/paseto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC *application.UserUseCase
}

func NewUserHandler(userUC *application.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// Call use-case
	user, err := h.userUC.GetUserByID(uint(id))
	if err != nil {
		// For example, gorm.ErrRecordNotFound => 404
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CheckIfUserIsCompany(c *gin.Context) {
	login := c.Query("login")
	isCompany, err := h.userUC.CheckIfUserIsCompany(login)
	if err != nil {
		if h.userUC.GetRepo().IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if user is company"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": isCompany})
}

func (h *UserHandler) CheckIfUserIsLogged(c *gin.Context) {
	login := c.Query("login")
	isLogged, err := h.userUC.CheckIfUserIsLogged(login)
	if err != nil {
		if h.userUC.GetRepo().IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if user is logged"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": isLogged})
}

func (h *UserHandler) GetUserByLogin(c *gin.Context) {
	login := c.Query("login")
	user, err := h.userUC.GetUserByLogin(login)
	if err != nil {
		if h.userUC.GetRepo().IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user by login"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserAndSubUsersByOwnerUsername(c *gin.Context) {

	claims, err := paseto.Paseto().ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var mainUserUsername string

	if claims.OwnerUsername == "" {
		mainUserUsername = claims.Username
	} else {
		mainUserUsername = claims.OwnerUsername
	}

	mainUser, subUsers, err := h.userUC.GetUserAndSubUsersByOwnerUsername(mainUserUsername)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user and subusers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": mainUser, "subusers": subUsers})
}
