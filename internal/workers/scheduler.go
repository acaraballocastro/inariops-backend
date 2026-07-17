package workers

import (
	"context"
	"log"
	"time"
)

type Scheduler struct {
	workers []Worker
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		workers: []Worker{},
	}
}

func (s *Scheduler) Register(worker Worker) {
	s.workers = append(s.workers, worker)
}

func (s *Scheduler) Start(ctx context.Context) {
	for _, worker := range s.workers {
		go s.startWorker(ctx, worker)
	}
}

func (s *Scheduler) startWorker(ctx context.Context, worker Worker) {
	log.Printf("[WORKER] %s started", worker.Name())

	ticker := time.NewTicker(worker.Interval())
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			log.Printf("[WORKER] %s stopped", worker.Name())
			return

		case <-ticker.C:
			worker.Run(ctx)

		}
	}
}
