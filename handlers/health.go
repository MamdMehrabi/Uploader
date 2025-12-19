package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/MamdMehrabi/Uploader/models"
)

func HealthHandler(botToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		botStatus := "missing"
		if botToken != "" {
			if len(botToken) >= 10 {
				botStatus = "configured"
			} else {
				botStatus = "invalid (too short)"
			}
		}
		c.JSON(http.StatusOK, models.HealthResponse{
			Status:   "ok",
			BotToken: botStatus,
		})
	}
}

func MaxFileSizeHandler(c *gin.Context) {
	maxSizeMB := c.MustGet("maxFileSizeMB").(int)
	c.JSON(200, models.MaxFileSizeResponse{
		MaxFileSizeMB:    maxSizeMB,
		MaxFileSizeBytes: c.MustGet("maxFileSizeBytes").(int64),
	})
}

