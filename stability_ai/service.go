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

// GenerateImage Generate multiple images (Concurrency)
func GenerateImage(engineId string, generateInput *GenerateInput) (Images, error) {
	reqURL := ApiHost + "/v1alpha/generation/" + engineId + "/text-to-image"
	var images Images

	var wg sync.WaitGroup
	var mutex sync.Mutex

	done := make(chan struct{})
	errCh := make(chan error)

	for i := 0; i < generateInput.Samples; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			stabilityApiPayload := generateInput.ToStabilityApiPayload()
			stabilityApiPayload.Samples = 1

			resp, err := client.R().
				SetHeader("Accept", "image/png").
				SetHeader("Authorization", apiKey).
				SetBody(stabilityApiPayload).
				Post(reqURL)

			if err != nil {
				errCh <- err
				return
			}

			if resp.StatusCode() != http.StatusOK {
				errCh <- fmt.Errorf("HTTP %d: %s", resp.StatusCode(), resp.Body())
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
				Msgf("Finished creating the %d st image", len(images))
		case err := <-errCh:
			log.Error().Caller().Msg(err.Error())
		case <-time.After(60 * time.Second):
			log.Error().Caller().Msg("Timeout exceeded while generating images")
		}
	}

	return images, nil
}

// Engines Lookup available engines
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

// Balance Get the current account's remaining credit
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

// calcCredit Calculating and logging credit consumption
func calcCredit(c *gin.Context, done chan struct{}, req *StabilityApiPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	before, _ := Balance()

	select {
	case <-done:
		after, _ := Balance()

		// 1,000 credit = $10
		cost := before.Credits - after.Credits

		summary := map[string]any{
			"2. Number of images":  req.Samples,
			"3. Steps":             req.Steps,
			"4. Credits consumed":  cost,
			"5. Dollar ($)":        cost * 0.01,
			"7. Remaining credits": after.Credits,
		}
		for _, prompt := range req.TextPrompts {
			summary["1. Prompt"] = prompt.Text
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
