package middleware

import "github.com/gin-gonic/gin"

func SetHeader(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "no-store")
	c.Next()
}

func AllowPreflight(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
}
