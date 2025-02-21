package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	count, err := h.store.GetVisitCount(shortCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "short link not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"visits": count})
}
