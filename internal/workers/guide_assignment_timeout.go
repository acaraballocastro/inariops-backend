package workers

import (
	"context"
	"time"

	"inariops/internal/domain"
	tourdaysapp "inariops/internal/modules/tours/application/tour_days"
	"inariops/internal/shared/logger"
)

type GuideAssignmentTimeoutWorker struct {
	service *tourdaysapp.Service
}

func NewGuideAssignmentTimeoutWorker(service *tourdaysapp.Service) *GuideAssignmentTimeoutWorker {
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

func (w *GuideAssignmentTimeoutWorker) Run(ctx context.Context) (domain.Result, error) {
	var result domain.Result

	logger.Debug("[WORKER] Checking expired guide assignments...")

	result, err := w.service.ProcessExpiredGuideAssignments(ctx)
	if err != nil {
		logger.Error("[WORKER] Failed to process expired guide assignments: %v", err)
	}

	logger.Info(
		"[WORKER] Assigned guide assignments: found=%d processed=%d failed=%d",
		result.Found,
		result.Processed,
		result.Failed,
	)

	return result, nil
}
