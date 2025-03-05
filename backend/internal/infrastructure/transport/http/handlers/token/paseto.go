package token

import (
	"app/internal/infrastructure/token/paseto"
	"app/pkg/errorsLib"
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
		ctx.JSON(errorsLib.HTTPStatusCode(err.Error()), gin.H{"error": err.Error()})
		return
	}

	if ctx.Query("code") != "123456" {
		ctx.JSON(errorsLib.HTTPStatusCode(errorsLib.ErrAccessDenied.Error()), gin.H{"error": errorsLib.ErrAccessDenied.Error()})
		return
	}

	ctx.JSON(http.StatusOK, claims)
}
