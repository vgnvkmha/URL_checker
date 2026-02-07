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

func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{
			Timeout: timeout, // запасной таймаут
		},
	}
}

func (c *HTTPChecker) Check(
	ctx context.Context,
	target dto.Targets,
) (dto.Checks, error) {

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.URL, nil)
	if err != nil {
		return dto.Checks{}, err
	}

	resp, err := c.client.Do(req)
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
