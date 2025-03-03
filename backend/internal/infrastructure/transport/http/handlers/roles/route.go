package roles

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	handler := NewRoleHandler(application.NewRoleUseCase(repositories.NewRoleRepository(), repositories.NewUserRepository()))

	// // Routes
	group := router.Group("/roles")
	{
		group.GET("/all", handler.GetAllRoles)                // Get all roles
		group.GET("by-id", handler.GetRoleByID)               // Get role by ID
		group.GET("/by-username", handler.GetRolesByUsername) // Get roles by username

		group.POST("assign", handler.AssignRolesToUser) // Assign roles to user
		group.POST("remove", handler.RemoveRolesOfUser) // Remove roles from user
	}
}
