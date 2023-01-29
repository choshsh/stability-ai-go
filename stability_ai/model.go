package stability_ai

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Image struct {
	ImageID     uuid.UUID            `dynamo:",hash" index:"Seq-ID-index,hash" json:"imageID"` // 이미지 고유 ID
	Time        time.Time            `dynamo:",range" json:"time"`
	Url         string               `json:"url"`         // 이미지 URL
	Keywords    []string             `json:"keywords"`    // 이미지의 키워드 (프롬프트에서 추출)
	RequestInfo *StabilityApiPayload `json:"requestInfo"` // 이미지 생성 요청값
}

// SetKeywords 프롬프트에서 키워드 추출
func (i *Image) SetKeywords() {
	if i.RequestInfo != nil {
		split := strings.Split(i.RequestInfo.TextPrompts[0].Text, ",")
		for _, s := range split {
			i.Keywords = append(i.Keywords, strings.TrimSpace(s))
		}
	}
}

// SetUrlPrefix S3 키 값에 접근 가능한 도메인을 추가
func (i *Image) SetUrlPrefix() {
	i.Url = strings.Join([]string{cloudfrontHost, i.Url}, "")
}