package handlers

import (
	"app/internal/application"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	providerUC *application.ProviderUseCase
}

func NewProviderHandler(providerUC *application.ProviderUseCase) *ProviderHandler {
	return &ProviderHandler{providerUC: providerUC}
}

func (h *ProviderHandler) GetAllProviders(c *gin.Context) {
	providers, err := h.providerUC.GetAllProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, providers)
}

func (h *ProviderHandler) GetProviderByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider ID"})
		return
	}

	provider, err := h.providerUC.GetProviderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, provider)
}
