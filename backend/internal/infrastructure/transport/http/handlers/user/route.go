package user

import (
	"app/internal/application"
	"app/internal/infrastructure/repositories"
	"app/internal/infrastructure/webhooks/verificaciones"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	handler := NewUserHandler(
		application.NewUserUseCase(
			repositories.NewUserRepository(),
			repositories.NewInternalCompanyRepository(),
			verificaciones.NewVerificacionesClient()))

	subUserHandler := NewSubUserHandler(
		application.NewSubUserUseCase(
			repositories.NewUserRepository(),
			repositories.NewRoleRepository(),
			verificaciones.NewVerificacionesClient()))
	// // Routes
	group := router.Group("/users")
	{
		// // Get all providers
		// group.GET("/all", handler.GetAllProviders)
		group.POST("/register", handler.RegisterCompanyUser) // Register company user
		group.POST("/subuser", subUserHandler.CreateSubUser) // Create subuser

		group.GET("/all", handler.GetUserAndSubUsersByOwnerUsername) // Get user and subusers by owner username
		group.GET("/by-id", handler.GetUserByID)                     // Get user by ID
		group.GET("/by-login", handler.GetUserByLogin)               // Get user by login
		group.GET("/is-company", handler.CheckIfUserIsCompany)       // Check if user is company
		group.GET("/is-logged", handler.CheckIfUserIsLogged)         // Check if user is logged

		group.POST("/activation", handler.ActivateDeactivateUser) // Activate/deactivate user
	}
}
