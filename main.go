package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"stability-ai-go/docs"
	"stability-ai-go/middleware"
	"stability-ai-go/stability_ai"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
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

	ginLambda = ginadapter.NewV2(r)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
