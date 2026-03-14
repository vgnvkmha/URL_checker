package writer

import (
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"context"

	"go.uber.org/zap"
)

type Writer struct {
	repo    checks.ICheckRepository
	results <-chan dto.Checks
	logger  *zap.SugaredLogger
}

func NewWriter(
	repo checks.ICheckRepository,
	results <-chan dto.Checks,
	logger *zap.SugaredLogger) *Writer {

	return &Writer{
		repo:    repo,
		results: results,
		logger: logger.With("module", "writer",
			"layer", "writer")}
}

func (w *Writer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-w.results:
			w.logger.Infow("Writting Info",
				"target_id", res.TargetId)
			if _, err := w.repo.Insert(ctx, res); err != nil {
				w.logger.Errorw("Writting failed",
					"error", err,
					"target_id", res.TargetId)
			}
		}
	}
}
