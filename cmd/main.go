package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"stability-ai-go/docs"
	"stability-ai-go/middleware"
	"stability-ai-go/stability_ai"
)

func main() {
	r := gin.Default()

	// Middleware
	r.Use(middleware.SetHeader)
	r.Use(middleware.AllowPreflight)

	// Router
	stability_ai.NewRouterV1(r)

	docs.SwaggerInfo.Title = "Stability AI - Demo"
	docs.SwaggerInfo.Description = "Contact. cho911115@gmail.com"
	docs.SwaggerInfo.Version = "1.0"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(99)),
	)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
