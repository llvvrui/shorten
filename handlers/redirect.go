package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")
	
	longUrl, err := h.store.GetLongUrl(shortCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "short link not found"})
		return
	}

	_ = h.store.IncrementVisit(shortCode)
	c.Redirect(http.StatusFound, longUrl)
}
