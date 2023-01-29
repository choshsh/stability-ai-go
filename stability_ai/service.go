package stability_ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"sync"
	"time"

	"net/http"
)

// GenerateImage n개의 이미지 생성 (Concurrency)
func GenerateImage(engineId string, generateInput *GenerateInput) ([]*Image, error) {
	reqUrl := ApiHost + "/v1alpha/generation/" + engineId + "/text-to-image"

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	var images []*Image

	for i := 0; i < generateInput.Samples; i++ {
		wg.Add(1)

		go func(ctx context.Context, wg *sync.WaitGroup) {
			done := make(chan struct{})

			stabilityApiPayload := generateInput.ToStabilityApiPayload()
			stabilityApiPayload.Samples = 1

			go func() {
				resp, err := client.R().
					SetHeader("Accept", "image/png").
					SetHeader("Authorization", apiKey).
					SetBody(stabilityApiPayload).
					Post(reqUrl)

				if err != nil {
					log.Error().Caller().Err(err).Msg("")
				}

				if resp.StatusCode() != http.StatusOK {
					log.Error().Caller().Msg(string(resp.Body()))
				}

				save, _ := SaveImage(resp.Body(), stabilityApiPayload)
				mutex.Lock()
				images = append(images, save)
				mutex.Unlock()

				done <- struct{}{}
			}()

			select {
			case <-done:
				log.Info().
					Str("Prompt", stabilityApiPayload.TextPrompts[0].Text).
					Str("Url", images[len(images)-1].Url).
					Msgf("%d번째 이미지 생성 완료", len(images))
			case <-ctx.Done():
				log.Error().Caller().Msg(ctx.Err().Error())
			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return images, nil
}

// Engines 사용 가능한 engine 조회
func Engines() (*EnginesResponse, error) {
	reqUrl := ApiHost + "/v1alpha/engines/list"

	enginesResponse := EnginesResponse{}

	resp, _ := client.R().
		SetHeader("Authorization", apiKey).
		SetResult(&enginesResponse).
		Get(reqUrl)

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	return &enginesResponse, nil
}

// Balance 현재 계정의 잔여 크레딧 조회
func Balance() (*BalanceResponse, error) {
	reqUrl := ApiHost + "/v1alpha/user/balance"

	balanceResponse := BalanceResponse{}

	resp, _ := client.R().
		SetHeader("Authorization", apiKey).
		SetResult(&balanceResponse).
		Get(reqUrl)

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	return &balanceResponse, nil
}

// calcCredit 크레딧 소모량 계산 및 로깅
func calcCredit(c *gin.Context, done chan struct{}, req *StabilityApiPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	before, _ := Balance()

	select {
	case <-done:
		after, _ := Balance()

		// 1,000 credit = $10
		// 환율: 1,250
		cost := before.Credits - after.Credits

		summary := map[string]any{
			"2. 이미지 개수":   req.Samples,
			"3. Steps":    req.Steps,
			"4. 크레딧 소모":   cost,
			"5. 달러 ($)":   cost * 0.01,
			"6. 원화 (₩)":   cost * 0.01 * 1250,
			"7. 잔여 크레딧: ": after.Credits,
		}
		for _, prompt := range req.TextPrompts {
			summary["1. 프롬프트"] = prompt.Text
		}

		marshal, _ := json.Marshal(summary)
		fmt.Println(string(marshal))

	case <-c.Done():
		if c.Err() != nil {
			log.Error().Caller().Err(c.Err()).Msg("")
		}
	case <-ctx.Done():
		close(done)
		log.Error().Caller().Err(ctx.Err()).Msg("")
	}
}
