package profile

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	handler := NewProfileHandler(application.NewProfileUseCase(repositories.NewUserRepository()))

	// // Routes
	group := router.Group("/users")
	{
		// // Get all providers
		// group.GET("/all", handler.GetAllProviders)

		group.POST("/profile/upload", handler.UpdateOwnProfile)       // Upload profile
		group.POST("/profile/by-username", handler.UpdateUserProfile) // Update user profile

	}
}
