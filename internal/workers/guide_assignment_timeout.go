package workers

import (
	"context"
	"log"
	"time"

	tourdays "inariops/internal/modules/tours/tour_days"
)

type GuideAssignmentTimeoutWorker struct {
	service *tourdays.Service
}

func NewGuideAssignmentTimeoutWorker(service *tourdays.Service) *GuideAssignmentTimeoutWorker {
	return &GuideAssignmentTimeoutWorker{
		service: service,
	}
}

func (w *GuideAssignmentTimeoutWorker) Name() string {
	return "guide-assignment-timeout"
}

func (w *GuideAssignmentTimeoutWorker) Interval() time.Duration {
	return time.Minute
}

func (w *GuideAssignmentTimeoutWorker) Run(ctx context.Context) {

	log.Println("[WORKER] Checking expired guide assignments...")

	err := w.service.ProcessExpiredGuideAssignments(ctx)
	if err != nil {
		log.Println(err)
	}
}
