package handlers

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"

	"github.com/gin-gonic/gin"
)

func ProviderRoutes(router *gin.Engine) {
	repo := repositories.NewProviderRepository()
	usecase := application.NewProviderUseCase(repo)
	handler := NewProviderHandler(usecase)

	// // Routes
	group := router.Group("/providers")
	{
		// Get all providers
		group.GET("/all", handler.GetAllProviders)

		// Get provider by ID
		group.GET("", handler.GetProviderByID)

	}
}
