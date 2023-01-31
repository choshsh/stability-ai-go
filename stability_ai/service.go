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
func GenerateImage(engineId string, generateInput *GenerateInput) (Images, error) {
	reqUrl := ApiHost + "/v1alpha/generation/" + engineId + "/text-to-image"
	var images Images

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	done := make(chan struct{})
	errCh := make(chan error)

	for i := 0; i < generateInput.Samples; i++ {
		wg.Add(1)
		stabilityApiPayload := generateInput.ToStabilityApiPayload()
		stabilityApiPayload.Samples = 1

		go func() {
			defer wg.Done()

			resp, err := client.R().
				SetHeader("Accept", "image/png").
				SetHeader("Authorization", apiKey).
				SetBody(stabilityApiPayload).
				Post(reqUrl)

			if err != nil {
				errCh <- err
				return
			}

			if resp.StatusCode() != http.StatusOK {
				errCh <- errors.New(string(resp.Body()))
				return
			}

			save, _ := SaveImage(resp.Body(), stabilityApiPayload)

			mutex.Lock()
			images = append(images, save)
			mutex.Unlock()

			done <- struct{}{}
		}()
	}

	go func() {
		wg.Wait()
		close(done)
		close(errCh)
	}()

	for i := 0; i < generateInput.Samples; i++ {
		select {
		case <-done:
			log.Info().
				Str("Prompt", generateInput.TextPrompt).
				Str("Url", images[len(images)-1].Url).
				Msgf("%d번째 이미지 생성 완료", len(images))
		case err := <-errCh:
			log.Error().Caller().Msg(err.Error())
		case <-time.After(60 * time.Second):
			log.Error().Caller().Msg("이미지 생성 에러 [timeout 60s]")
		}
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
