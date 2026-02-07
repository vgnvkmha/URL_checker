package writer

import (
	"URL_checker/internal/repo/checks"
	"URL_checker/internal/repo/dto"
	"context"
	"fmt"
	"log"
)

type Writer struct {
	repo    checks.ICheckRepository
	results <-chan dto.Checks
}

func NewWriter(
	repo checks.ICheckRepository,
	results <-chan dto.Checks,
) *Writer {
	return &Writer{repo, results}
}

func (w *Writer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case res := <-w.results:
			fmt.Printf("WRITER INSERT: target_id=%d", res.TargetId)
			if _, err := w.repo.Insert(ctx, res); err != nil {
				log.Println("writer error:", err)
			}
		}
	}
}
