package api

import (
	"time"

	"golang.org/x/time/rate"
)

// NewLimiter guards the paid model path: burst calls at once, then one per interval.
func NewLimiter(interval time.Duration, burst int) *rate.Limiter {
	return rate.NewLimiter(rate.Every(interval), burst)
}
