package stability_ai

import "github.com/gin-gonic/gin"

func NewRouterV1(r *gin.Engine) *gin.RouterGroup {
	v1 := r.Group("/v1")
	{
		v1.GET("/image", FindImageManyCtrl)
		v1.GET("/image/:id", FindImageByIdCtrl)
		v1.POST("/generate/:engineId", GenerateImageCtrl)
		v1.GET("/balance", BalanceCtrl)
	}
	return v1
}
