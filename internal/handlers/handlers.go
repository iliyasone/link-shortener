package handlers

import (
	"net/http"

	"link-shortener/internal/service"
	"link-shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	OriginalURL string `json:"originalURL"`
}

type ShortenResponse struct {
	ShortURL string `json:"shortURL,omitempty"`
	Error    string `json:"error,omitempty"`
}

func PostShortenURL(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON input"})
			return
		}

		if req.OriginalURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "originalURL is required"})
			return
		}

		shortURL, err := service.SaveURL(store, req.OriginalURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, ShortenResponse{ShortURL: shortURL})
	}
}

func GetOriginalURL(store storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		if shortURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "shortURL parameter is required"})
			return
		}

		originalURL, err := store.Get(shortURL)
		if err != nil {
			// e.g. storage.ErrURLNotFound
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"originalURL": originalURL})
	}
}
