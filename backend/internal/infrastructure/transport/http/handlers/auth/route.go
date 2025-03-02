package auth

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"
	"app/internal/infrastructure/webhooks/verificaciones"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	handler := NewAuthHandler(application.NewAuthUseCase(repositories.NewUserRepository(), verificaciones.NewVerificacionesClient()))
	// // Routes
	group := router.Group("/auth")
	{
		// // Get all providers
		// group.GET("/all", handler.GetAllProviders)

		// Get user by ID
		group.POST("/login", handler.Login)

	}
}
