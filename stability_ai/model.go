package stability_ai

import (
	"github.com/google/uuid"
	"sort"
	"strings"
	"time"
)

const (
	ModelImagePK           = "ImageID"
	ModelImagePutCondition = "attribute_not_exists(" + ModelImagePK + ")"
)

type Image struct {
	ImageID     uuid.UUID            `dynamo:",hash" index:"Seq-ID-index,hash" json:"imageID"` // 이미지 고유 ID
	Time        time.Time            `dynamo:",range" json:"time"`
	Url         string               `json:"url"`         // 이미지 URL
	Keywords    []string             `json:"keywords"`    // 이미지의 키워드 (프롬프트에서 추출)
	RequestInfo *StabilityApiPayload `json:"requestInfo"` // 이미지 생성 요청값
}

// SetKeywords Extract keywords from prompts
func (i *Image) SetKeywords() *Image {
	if i.RequestInfo != nil {
		split := strings.Split(i.RequestInfo.TextPrompts[0].Text, ",")
		for _, s := range split {
			i.Keywords = append(i.Keywords, strings.TrimSpace(s))
		}
	}

	return i
}

// SetUrlPrefix Add accessible domain to S3 key values
func (i *Image) SetUrlPrefix() *Image {
	i.Url = strings.Join([]string{cloudfrontHost, i.Url}, "")

	return i
}

type Images []*Image

// SortTimeDesc Sort Descending by Time
func (is Images) SortTimeDesc() Images {
	sort.Slice(is, func(i, j int) bool {
		return is[i].Time.After(is[j].Time)
	})

	return is
}

// SetUrlPrefix Add accessible domain to S3 key values
func (is Images) SetUrlPrefix() Images {
	for _, image := range is {
		image.Url = strings.Join([]string{cloudfrontHost, image.Url}, "")
	}

	return is
}
