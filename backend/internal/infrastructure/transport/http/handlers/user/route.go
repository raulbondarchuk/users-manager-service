package user

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	handler := NewUserHandler(application.NewUserUseCase(repositories.NewUserRepository()))
	subUserHandler := NewSubUserHandler(application.NewSubUserUseCase(repositories.NewUserRepository(), repositories.NewRoleRepository()))
	// // Routes
	group := router.Group("/users")
	{
		// // Get all providers
		// group.GET("/all", handler.GetAllProviders)

		// Get user by ID
		group.GET("", handler.GetUserByID)

		// Create subuser
		group.POST("/subuser", subUserHandler.CreateSubUser)
	}
}
