package checker

import (
	"URL_checker/internal/repo/dto"
	"context"
	"log"
	"net/http"
	"time"
)

type Checker interface {
	Check(ctx context.Context, target dto.Targets) (dto.Checks, error)
}

type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: NewHTTPClient(),
	}
}

func NewHTTPClient() *http.Client {
	tr := &http.Transport{
		ForceAttemptHTTP2: true, // как браузер
	}

	return &http.Client{
		Timeout:   15 * time.Second,
		Transport: tr,
	}
}

func (c *HTTPChecker) Check(
	ctx context.Context,
	target dto.Targets,
) (dto.Checks, error) {

	start := time.Now()

	req, _ := http.NewRequest(http.MethodGet, target.URL, nil)

	req.Header.Set("User-Agent", "curl/8.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start)
	log.Printf(
		"CHECK done target=%d ok=%v",
		target.ID,
		target.Active,
	)

	if err != nil {
		return dto.Checks{
			TargetId:  target.ID,
			CheckedAt: time.Now(),
			OK:        false,
			Error:     err.Error(),
			LatencyMs: latency.Milliseconds(),
		}, nil
	}
	defer resp.Body.Close()

	return dto.Checks{
		TargetId:   target.ID,
		CheckedAt:  time.Now(),
		OK:         resp.StatusCode < 500,
		StatusCode: uint8(resp.StatusCode),
		LatencyMs:  latency.Milliseconds(),
	}, nil
}
