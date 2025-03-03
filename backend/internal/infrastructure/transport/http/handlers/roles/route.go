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
		group.GET("/all", handler.GetAllRoles)                       // Get all roles
		group.GET("", handler.GetRoleByID)                           // Get role by ID
		group.GET("/username/:username", handler.GetRolesByUsername) // Get roles by username

	}
}
