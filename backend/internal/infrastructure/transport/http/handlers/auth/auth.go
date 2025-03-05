package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/internal/application"
	"app/pkg/errorsLib"
)

type AuthHandler struct {
	authUC *application.AuthUseCase
}

func NewAuthHandler(uc *application.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: uc}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	user, err := h.authUC.Login(req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !user.Active {
		c.JSON(errorsLib.HTTPStatusCode(errorsLib.ErrAccessDenied.Error()), gin.H{"error": errorsLib.ErrAccessDenied.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+user.AccessToken)
	c.Header("Refresh", *user.Refresh)

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) RefreshPairTokens(c *gin.Context) {
	refreshTokenReq := c.Query("refresh")

	accessToken, refreshToken, err := h.authUC.RefreshPairTokens(refreshTokenReq)
	if err != nil {
		if err.Error() == errorsLib.ErrAccessDenied.Error() {
			c.JSON(http.StatusLocked, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"access": accessToken, "refresh": refreshToken})
}
