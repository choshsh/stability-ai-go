package stability_ai

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GenerateImageCtrl
// @Summary 이미지 생성하기
// @Tags image
// @Accept json
// @Produce json
// @Param engineId path string true "엔진 ID" Enums(stable-diffusion-512-v2-1, stable-diffusion-512-v2-0, stable-diffusion-v1-5)
// @Param GenerateInput body GenerateInput true "상세는 아래 Model 클릭해서 정보 확인"
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

	// 이미지 생성
	result, err := GenerateImage(engineId, &generateInput)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NewBaseErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
	done <- struct{}{}
}

// FindImageByIdCtrl
// @Summary ID로 이미지 조회
// @Tags image
// @Accept json
// @Produce json
// @Param id path string true "image의 uuid"
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
// @Summary 이미지 리스트 조회
// @Tags image
// @Accept json
// @Produce json
// @Param keyword query string false "(Optional) keyword로 조건 검색"
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
// @Summary stability.ai 잔여 크레딧 조회
// @Tags credit
// @Accept json
// @Produce json
// @Success 200 {object} BalanceResponse
// @Failure 400 {object} BaseErrorResponse
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
