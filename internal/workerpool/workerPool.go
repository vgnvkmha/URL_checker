package workerpool

import (
	"URL_checker/internal/checker"
	"URL_checker/internal/repo/dto"
	"context"
	"log"
	"time"
)

type WorkerPool struct {
	checker checker.Checker
	workers int
	results chan<- dto.Checks // <- пишем результаты в канал
}

func New(checker checker.Checker, results chan<- dto.Checks, workers int) *WorkerPool {
	return &WorkerPool{
		checker: checker,
		workers: workers,
		results: results,
	}
}

func (wp *WorkerPool) Start(
	ctx context.Context,
	queue <-chan dto.Targets,
) {
	//fan-out
	for i := 0; i < wp.workers; i++ {
		go wp.worker(ctx, i, queue)
	}
}

func (wp *WorkerPool) worker(ctx context.Context, id int, queue <-chan dto.Targets) {
	for {
		select {
		case <-ctx.Done():
			return
		case target := <-queue:
			checkCtx, cancel := context.WithTimeout(
				ctx,
				time.Duration(target.TimeoutMS)*time.Millisecond,
			)
			result, err := wp.checker.Check(checkCtx, target)
			cancel()

			if err != nil {
				log.Println("check error:", err)
				continue
			}

			log.Printf("CHECK done target=%d url=%s ok=%v status=%d",
				target.ID, target.URL, result.OK, result.StatusCode)

			// ВАЖНО: отправляем результат в writer через канал
			select {
			case wp.results <- result:
			case <-ctx.Done():
				return
			}
		}
	}
}
