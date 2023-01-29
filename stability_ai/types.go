package stability_ai

import "strings"

type GenerateInput struct {
	CfgScale           int     `json:"cfg_scale" binding:"required,min=0,max=35" default:"7"` // 확산 프로세스가 프롬프트 텍스트에 얼마나 엄격하게 적용되는지 (값이 높을수록 이미지가 프롬프트에 더 가깝게 유지됨)
	ClipGuidancePreset string  `json:"clip_guidance_preset" binding:"required,oneof=FAST_BLUE FAST_GREEN NONE SIMPLE SLOW SLOWER SLOWEST"`
	Height             int     `json:"height" binding:"required,min=512,max=2048" default:"512"`
	Width              int     `json:"width" binding:"required,min=512,max=2048" default:"512"`
	Samples            int     `json:"samples" binding:"required,min=1,max=10" default:"1"`  // 생성할 이미지 수
	Steps              int     `json:"steps" binding:"required,min=10,max=150" default:"40"` // 실행할 확산 단계 수. 이미지 품질에 영향 큼. 비용에도 영향 큼.
	TextPrompt         string  `json:"text_prompts" binding:"required"`
	Style              *string `json:"style,omitempty" default:""` // 개발 중. 생략해도 됨
}

func (gi *GenerateInput) Preprocessing() {
	if gi.Style != nil {
		switch *gi.Style {
		case StylePhoto:
			gi.TextPrompt = strings.Join([]string{StylePhoto, gi.TextPrompt}, "")
		case StyleRealistic:
			gi.TextPrompt = strings.Join([]string{gi.TextPrompt, StyleRealistic}, "")
		}
	}
}

func (gi *GenerateInput) ToStabilityApiPayload() *StabilityApiPayload {
	var textPrompts []*TextPrompt
	textPrompts = append(textPrompts, &TextPrompt{
		Text:   gi.TextPrompt,
		Weight: 1,
	})

	return &StabilityApiPayload{
		CfgScale:           gi.CfgScale,
		ClipGuidancePreset: gi.ClipGuidancePreset,
		Height:             gi.Height,
		Width:              gi.Width,
		Samples:            gi.Samples,
		Steps:              gi.Steps,
		TextPrompts:        textPrompts,
	}
}

// StabilityApiPayload Stability AI API - 이미지 생성 파라미터
type StabilityApiPayload struct {
	CfgScale           int           `json:"cfg_scale" binding:"required"`
	ClipGuidancePreset string        `json:"clip_guidance_preset"`
	Height             int           `json:"height"`
	Width              int           `json:"width"`
	Samples            int           `json:"samples" binding:"required"`
	Steps              int           `json:"steps"`
	TextPrompts        []*TextPrompt `json:"text_prompts" binding:"required"`
}

type TextPrompt struct {
	Text   string `json:"text"`
	Weight int    `json:"weight" default:"1"`
}

type EnginesResponse struct {
	Engines []struct {
		Description string `json:"description"`
		Id          string `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
	} `json:"engines"`
}

// StabilityApiResponse Stability AI API - 에러 reponse
type StabilityApiResponse struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	Message   string `json:"message"`
	Temporary bool   `json:"temporary"`
	Timeout   bool   `json:"timeout"`
	Fault     bool   `json:"fault"`
}

type BalanceResponse struct {
	Credits float64 `json:"credits"`
}

type BaseErrorResponse struct {
	Message string `json:"message"`
}

func NewBaseErrorResponse(msg string) *BaseErrorResponse {
	return &BaseErrorResponse{Message: msg}
}
