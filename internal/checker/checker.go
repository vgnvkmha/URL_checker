package checker

import (
	"URL_checker/internal/repo/dto"
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Checker interface {
	Check(ctx context.Context, target dto.Targets) (dto.Checks, error)
}

type HTTPChecker struct {
	client *http.Client
	logger *zap.SugaredLogger
}

func NewHTTPChecker(logger *zap.SugaredLogger) *HTTPChecker {
	return &HTTPChecker{
		client: NewHTTPClient(),
		logger: logger.With("module", "checker"),
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

	const operation = "Check"
	start := time.Now()

	req, _ := http.NewRequest(http.MethodGet, target.URL, nil)

	req.Header.Set("User-Agent", "curl/8.0") //TODO: Убрать захардкоженные хеддеры
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start)

	if err != nil {
		const red = "\033[31m"
		const reset = "\033[0m"
		c.logger.Infow(red+"Failed"+reset,
			"operation", operation,
			"target_id", target.ID,
		)
		return dto.Checks{
			TargetId:  target.ID,
			CheckedAt: time.Now(),
			OK:        false,
			Error:     err.Error(),
			LatencyMs: latency.Milliseconds(),
		}, nil
	}

	defer resp.Body.Close()

	c.logger.Infow("Succesful",
		"operation", operation,
		"target_id", target.ID,
	)

	return dto.Checks{
		TargetId:   target.ID,
		CheckedAt:  time.Now(),
		OK:         resp.StatusCode < 500,
		StatusCode: uint8(resp.StatusCode),
		LatencyMs:  latency.Milliseconds(),
	}, nil
}
