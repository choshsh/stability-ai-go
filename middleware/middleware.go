package middleware

import "github.com/gin-gonic/gin"

func SetHeader(c *gin.Context) {
	// lambda url에서 cors 설정을 할 경우 중복된 헤더 값으로 들어가서 클라이언트 에러 발생
	//c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "no-store")
	c.Next()
}

func AllowPreflight(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}
}
