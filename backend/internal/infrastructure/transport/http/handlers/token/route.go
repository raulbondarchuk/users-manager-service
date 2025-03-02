package token

import (
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {

	handler := NewPasetoHandler()

	// // Routes
	group := router.Group("/token")
	{
		// Get all providers
		group.GET("/decode", handler.DecodePasetoToken)
	}
}
