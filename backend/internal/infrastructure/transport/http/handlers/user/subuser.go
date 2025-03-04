package user

import (
	"net/http"

	"app/internal/application"
	"app/internal/infrastructure/token/paseto"

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

	subUser, err := h.subUserUseCase.CreateSubUser(claims.Username, req.Username, req.Password, req.Roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subUser)
}
