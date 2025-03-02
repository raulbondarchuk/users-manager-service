package http

import (
	"app/internal/infrastructure/transport/http/handlers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	handlers.ProviderRoutes(router)
	printRoutes(router)
}

// func providerRoutes(router *gin.Engine) {
// 	providerRepo := repositories.NewProviderRepository()
// 	providerUC := application.NewProviderUseCase(providerRepo)
// 	providerHandler := handlers.NewProviderHandler(providerUC)
// 	router.GET("/providers", providerHandler.GetAllProviders)
// }
