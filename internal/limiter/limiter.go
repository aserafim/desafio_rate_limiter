package limiter

import (
	"context"
)

type LimiterStrategy interface {
	Allow(ctx context.Context, key string, limit, blockDuration int) (bool, error)
}
