package workers

import (
	"context"
	"inariops/internal/domain"
	"time"
)

type Worker interface {
	Name() string
	Interval() time.Duration
	Run(ctx context.Context) (domain.Result, error)
}
