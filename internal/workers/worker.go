package workers

import (
	"context"
	"time"
)

type Worker interface {
	Name() string
	Interval() time.Duration
	Run(ctx context.Context)
}
