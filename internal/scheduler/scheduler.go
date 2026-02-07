package scheduler

import (
	"URL_checker/internal/repo/dto"
	"URL_checker/internal/repo/target"
	"context"
	"log"
	"time"
)

type Scheduler struct {
	queue      chan<- dto.Targets
	targetRepo target.ITargetRepo
}

func New(
	queue chan dto.Targets,
	repo target.ITargetRepo,
) *Scheduler {
	return &Scheduler{
		queue:      queue,
		targetRepo: repo,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	schedule := make(map[uint64]time.Time)

	// for _, t := range targets {
	// 	schedule[t.ID] = struct{}{}
	// }

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			targets, err := s.targetRepo.ListActive(ctx)
			if err != nil {
				log.Println("scheduler:", err)
				continue
			}

			now := time.Now()

			for _, t := range targets {
				log.Printf(
					"SCHEDULE target=%d next=%s active=%v",
					t.ID,
					schedule[t.ID],
					t.Active,
				)
				next, exists := schedule[t.ID]
				if !exists || now.After(next) {
					select {
					case s.queue <- t:
						schedule[t.ID] = now.Add(time.Duration(t.IntervalSec) * time.Second)
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}
}
