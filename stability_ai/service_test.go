package stability_ai

import (
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
	"stability-ai-go/db"
	"testing"
)

func TestGenerateSingle(t *testing.T) {
	db.Init("local")

	prompts := `cat`

	input := GenerateInput{
		CfgScale:           7, // 프롬프트와의 근사값
		ClipGuidancePreset: "FAST_BLUE",
		Height:             512,
		Width:              512,
		Samples:            1,  // 몇 개의 이미지 만들지?
		Steps:              80, // 노이즈 제거 + 얼마나 디테일한
		TextPrompt:         prompts,
	}

	_, err := GenerateImage(ENGINE_512_V2_1, &input)
	assert.NoError(t, err)
}

func TestEngines(t *testing.T) {
	db.Init("local")

	result, err := Engines()
	fmt.Printf("Total count: %d\n\n", len(result.Engines))
	for _, engine := range result.Engines {
		pp.Printf("Name: %s, ID: %s\n", engine.Name, engine.Id)
	}
	assert.NoError(t, err)
}

func TestBalance(t *testing.T) {
	db.Init("local")

	result, err := Balance()

	fmt.Printf("Credit: %f\n\n", result.Credits)
	assert.NoError(t, err)
}

func TestFindMany(t *testing.T) {
	many, err := FindImageMany(nil)
	assert.NoError(t, err)

	pp.Println("size: ", len(many))
}
