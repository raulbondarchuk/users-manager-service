package user

import (
	"app/internal/application"
	"app/internal/infrastructure/token/paseto"
	"app/pkg/errorsLib"
	"net/http"
	"strconv"
	"strings"

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
		c.JSON(errorsLib.HTTPStatusCode(err.Error()), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": mainUser, "subusers": subUsers})
}

func (h *UserHandler) ActivateDeactivateUser(c *gin.Context) {
	username := c.Query("username")
	active := c.Query("active")

	activeBool, err := strconv.ParseBool(active)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid active value"})
		return
	}

	err = h.userUC.ActivateDeactivateUser(username, activeBool)
	if err != nil {
		if h.userUC.GetRepo().IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate/deactivate user"})
		}
		return
	}
	if activeBool {
		c.JSON(http.StatusOK, gin.H{"message": "User activated"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User deactivated"})
	}
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CompanyName string `json:"companyName" binding:"required"`
}

func (h *UserHandler) RegisterCompanyUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUC.RegisterCompanyUser(req.Username, req.Password, req.CompanyName)
	if err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
