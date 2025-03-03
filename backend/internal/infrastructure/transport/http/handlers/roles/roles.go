package roles

import (
	"app/internal/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleHandler - HTTP handler for roles
type RoleHandler struct {
	RoleUseCase *application.RoleUseCase
}

// NewRoleHandler - constructor for role handler
func NewRoleHandler(roleUseCase *application.RoleUseCase) *RoleHandler {
	return &RoleHandler{RoleUseCase: roleUseCase}
}

// GetAllRoles - handler for getting all roles
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.RoleUseCase.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GetRoleByID - handler for getting role by ID
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	role, err := h.RoleUseCase.GetRoleByID(uint(roleID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// GetRolesByUsername - handler for getting roles by username
func (h *RoleHandler) GetRolesByUsername(c *gin.Context) {
	username := c.Query("username")

	roles, err := h.RoleUseCase.GetRolesByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}
