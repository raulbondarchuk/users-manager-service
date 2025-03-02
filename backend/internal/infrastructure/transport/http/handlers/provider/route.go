package provider

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	// repo := repositories.NewProviderRepository()
	// usecase := application.NewProviderUseCase(repo)
	// handler := NewProviderHandler(usecase)
	handler := NewProviderHandler(application.NewProviderUseCase(repositories.NewProviderRepository()))

	// // Routes
	group := router.Group("/providers")
	{
		// Get all providers
		group.GET("/all", handler.GetAllProviders)

		// Get provider by ID
		group.GET("", handler.GetProviderByID)

	}
}
