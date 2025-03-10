package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"app/internal/application"
	"app/internal/infrastructure/token/paseto"
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

type forgotPasswordRequest struct {
	Username string `json:"username" binding:"required"`
	Link     string `json:"link" binding:"required"`
	Subject  string `json:"subject" binding:"required"`
	Body     string `json:"body" binding:"required"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {

	var req forgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	link, err := h.authUC.ForgotPassword(req.Username, req.Link, req.Subject, req.Body)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"link": link})
}

type resetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {

	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	claims, err := paseto.Paseto().ValidateToken(c.Query("token"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if claims.Roles != "recover" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	err = h.authUC.ResetPassword(claims.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
