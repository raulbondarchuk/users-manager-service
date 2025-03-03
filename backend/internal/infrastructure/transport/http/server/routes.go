package http

import (
	"app/internal/infrastructure/transport/http/handlers/auth"
	"app/internal/infrastructure/transport/http/handlers/provider"
	"app/internal/infrastructure/transport/http/handlers/roles"
	"app/internal/infrastructure/transport/http/handlers/token"
	"app/internal/infrastructure/transport/http/handlers/user"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	provider.Routes(router)
	user.Routes(router)
	auth.Routes(router)
	token.Routes(router)
	roles.Routes(router)

	printRoutes(router)
}
