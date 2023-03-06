package stability_ai

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GenerateImageCtrl
// @Summary Generate image
// @Tags image
// @Accept json
// @Produce json
// @Param engineId path string true "Engine ID" Enums(stable-diffusion-512-v2-1, stable-diffusion-512-v2-0, stable-diffusion-v1-5)
// @Param GenerateInput body GenerateInput true "Click on a model below for more information"
// @Success 200 {array} Image
// @Failure 400 {object} BaseErrorResponse
// @Failure 500 {object} BaseErrorResponse
// @Router /v1/generate/{engineId} [post]
func GenerateImageCtrl(c *gin.Context) {
	generateInput := GenerateInput{}
	if err := c.ShouldBindJSON(&generateInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseErrorResponse(err.Error()))
		return
	}

	generateInput.Preprocessing()
	engineId := c.Param("engineId")

	done := make(chan struct{})
	go calcCredit(c, done, generateInput.ToStabilityApiPayload())

	// Generate Image
	result, err := GenerateImage(engineId, &generateInput)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewBaseErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
	done <- struct{}{}
}

// FindImageByIdCtrl
// @Summary Find images by ID
// @Tags image
// @Accept json
// @Produce json
// @Param id path string true "uuid of the image"
// @Success 200 {object} Image
// @Failure 404 {object} BaseErrorResponse
// @Failure 500 {object} BaseErrorResponse
// @Router /v1/image/{id} [get]
func FindImageByIdCtrl(c *gin.Context) {
	id := c.Param("id")

	result, err := FindImageById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewBaseErrorResponse(err.Error()))
		return
	}
	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, NewBaseErrorResponse("item does not exist"))
		return
	}

	c.JSON(http.StatusOK, result)
}

// FindImageManyCtrl
// @Summary Get a list of images
// @Tags image
// @Accept json
// @Produce json
// @Param keyword query string false "(Optional) Search for conditions by keyword"
// @Success 200 {array} Image
// @Failure 404 {object} BaseErrorResponse
// @Failure 500 {object} BaseErrorResponse
// @Router /v1/image [get]
func FindImageManyCtrl(c *gin.Context) {
	var result Images
	var err error

	keyword := c.Query("keyword")

	result, err = FindImageMany(&FindManyInput{Keyword: keyword})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if result == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, NewBaseErrorResponse("item does not exist"))
		return
	}

	c.JSON(http.StatusOK, result)
}

// BalanceCtrl
// @Summary Get remaining credit
// @Tags credit
// @Produce json
// @Success 200 {object} BalanceResponse
// @Failure 500 {object} BaseErrorResponse
// @Router /v1/balance [get]
func BalanceCtrl(c *gin.Context) {
	result, err := Balance()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
