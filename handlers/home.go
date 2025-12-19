package handlers

import (
	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.File("./public/index.html")
}

