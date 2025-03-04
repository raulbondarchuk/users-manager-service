package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/internal/application"
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
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	c.Header("Authorization", "Bearer "+user.AccessToken)
	c.Header("Refresh", *user.Refresh)

	c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) RefreshPairTokens(c *gin.Context) {
	accessTokenExpiredReq := c.Query("access")
	refreshTokenReq := c.Query("refresh")

	accessToken, refreshToken, err := h.authUC.RefreshPairTokens(accessTokenExpiredReq, refreshTokenReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": accessToken, "refresh": refreshToken})
}
