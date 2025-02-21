package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/llvvrui/shortener/storage"
)

type ShortenRequest struct {
	LongUrl      string `json:"longUrl" binding:"required,url"`
	CustomSuffix string `json:"customSuffix,omitempty"`
	Expiration   int64  `json:"expiration,omitempty"`
}

type Handler struct {
	store storage.Store
}

func NewHandler(store storage.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Shorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiration time.Duration
	if req.Expiration > 0 {
		expiration = time.Duration(req.Expiration) * time.Second
	}

	shortCode, err := h.store.SaveShortUrl(req.CustomSuffix, req.LongUrl, expiration)
	if err != nil {
		if errors.Is(err, storage.ErrShortCodeExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "custom suffix already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create short link"})
		return
	}

	baseURL := "http://" + c.Request.Host + "/"
	c.JSON(http.StatusOK, gin.H{"shortUrl": baseURL + shortCode})
}