package user

import (
	"app/internal/application"
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
