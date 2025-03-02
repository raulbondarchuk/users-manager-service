package token

import (
	"app/internal/infrastructure/token/paseto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasetoHandler struct {
}

func NewPasetoHandler() *PasetoHandler {
	return &PasetoHandler{}
}

func (h *PasetoHandler) DecodePasetoToken(ctx *gin.Context) {

	claims, err := paseto.Paseto().ValidateToken(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ctx.Query("code") != "123456" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied"})
		return
	}

	ctx.JSON(http.StatusOK, claims)
}
