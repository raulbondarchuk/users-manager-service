package roles

import (
	"app/internal/application"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

// AssignRolesToUser - handler for assigning roles to user
func (h *RoleHandler) AssignRolesToUser(c *gin.Context) {
	// Parse request body
	var req struct {
		Username string `json:"username" binding:"required"`
		Roles    string `json:"roles" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call usecase
	if err := h.RoleUseCase.AssignRolesToUser(req.Username, req.Roles); err != nil {
		// Check if it's a "user not found" error
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Roles assigned to user %s successfully", req.Username),
	})
}
