package middleware

import "github.com/gin-gonic/gin"

func FileSizeLimit(maxFileSizeMB int) gin.HandlerFunc {
	maxFileSizeBytes := int64(maxFileSizeMB) * 1024 * 1024
	return func(c *gin.Context) {
		c.Set("maxFileSizeBytes", maxFileSizeBytes)
		c.Set("maxFileSizeMB", maxFileSizeMB)
		c.Next()
	}
}

