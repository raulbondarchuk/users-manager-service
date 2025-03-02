package http

import (
	"app/internal/infrastructure/transport/http/handlers/auth"
	"app/internal/infrastructure/transport/http/handlers/provider"
	"app/internal/infrastructure/transport/http/handlers/user"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	provider.Routes(router)
	user.Routes(router)
	auth.Routes(router)
	printRoutes(router)
}
