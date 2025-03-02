package http

import (
	provider "app/internal/infrastructure/transport/http/handlers/provider"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	provider.Routes(router)
	printRoutes(router)
}
