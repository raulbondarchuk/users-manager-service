package user

import (
	"net/http"
	"net/url"
	"strings"

	"app/internal/application"
	"app/internal/infrastructure/token/paseto"
	"app/pkg/config"
	"app/pkg/random"

	"github.com/gin-gonic/gin"
)

type SubUserHandler struct {
	subUserUseCase *application.SubUserUseCase
}

func NewSubUserHandler(subUserUseCase *application.SubUserUseCase) *SubUserHandler {
	return &SubUserHandler{
		subUserUseCase: subUserUseCase,
	}
}

type CreateSubUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Roles    string `json:"roles"`
	Email    string `json:"email"`
}

func (h *SubUserHandler) CreateSubUser(c *gin.Context) {

	claims, err := paseto.Paseto().ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var req CreateSubUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If middleware password is not correct, generate random password
	if config.ENV().MIDDLEWARE_PASSWORD == c.GetHeader("X-Middleware-Password") {
		req.Email = ""
	} else {
		req.Password, _ = random.GenerateRandomPassword()
	}

	// Create subuser
	subUser, err := h.subUserUseCase.CreateSubUser(claims.Username, req.Username, req.Password, req.Roles, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, subUser)
}

func (h *SubUserHandler) DeleteSubuser(c *gin.Context) {

	claims, err := paseto.Paseto().ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username of user is required"})
		return
	}

	// Decode the username if it is URL-encoded
	decodedUsername, err := url.QueryUnescape(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username format"})
		return
	}

	status, err := h.subUserUseCase.DeleteSubuser(decodedUsername, uint(claims.CompanyID))
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{"message": "Subuser deleted successfully"})
}
