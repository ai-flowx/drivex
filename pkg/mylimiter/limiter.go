package mylimiter

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"

	"github.com/ai-flowx/drivex/pkg/mycomdef"
)

var (
	limiterMap = make(map[string]*Limiter)
	mapMutex   sync.RWMutex
)

type Limiter struct {
	QPSLimiter         *rate.Limiter
	QPMLimiter         *SlidingWindowLimiter
	ConcurrencyLimiter *semaphore.Weighted
}

type SlidingWindowLimiter struct {
	mu          sync.Mutex
	maxRequests int
	interval    time.Duration
	requests    []time.Time
}

func NewSlidingWindowLimiter(qpm int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		maxRequests: qpm,
		interval:    time.Minute,
		requests:    make([]time.Time, 0, qpm),
	}
}

func (l *SlidingWindowLimiter) Allow() bool {
	now := time.Now()
	windowStart := now.Add(-l.interval)

	l.mu.Lock()
	defer l.mu.Unlock()

	i := 0

	for ; i < len(l.requests) && l.requests[i].Before(windowStart); i++ {
	}

	l.requests = l.requests[i:]

	if len(l.requests) < l.maxRequests {
		l.requests = append(l.requests, now)
		return true
	}

	return false
}

func (l *SlidingWindowLimiter) Wait(ctx context.Context) error {
	waitTime := 10 * time.Millisecond

	for {
		if l.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			l.mu.Lock()
			if len(l.requests) > 0 {
				nextAllowedTime := l.requests[0].Add(l.interval)
				timeUntilNextAllowed := time.Until(nextAllowedTime)
				if timeUntilNextAllowed < waitTime {
					waitTime = timeUntilNextAllowed
				} else {
					waitTime *= 2
					if waitTime > time.Second {
						waitTime = time.Second
					}
				}
			}
			l.mu.Unlock()
		}
	}
}

func NewLimiter(limitType string, limitn float64) *Limiter {
	lim := &Limiter{}

	switch limitType {
	case mycomdef.KeynameQps:
		lim.QPSLimiter = rate.NewLimiter(rate.Limit(limitn), int(limitn))
	case mycomdef.KeynameQpm, mycomdef.KeynameRpm:
		lim.QPMLimiter = NewSlidingWindowLimiter(int(limitn))
	case mycomdef.KeynameConcurrency:
		lim.ConcurrencyLimiter = semaphore.NewWeighted(int64(limitn))
	default:
	}

	return lim
}

func (l *Limiter) Wait(ctx context.Context) error {
	if l.QPSLimiter != nil {
		return l.QPSLimiter.Wait(ctx)
	}

	if l.QPMLimiter != nil {
		return l.QPMLimiter.Wait(ctx)
	}

	return nil
}

func (l *Limiter) Acquire(ctx context.Context) error {
	if l.ConcurrencyLimiter != nil {
		return l.ConcurrencyLimiter.Acquire(ctx, 1)
	}

	return nil
}

func (l *Limiter) Release() {
	if l.ConcurrencyLimiter != nil {
		l.ConcurrencyLimiter.Release(1)
	}
}

func GetLimiter(key, limitType string, limitn float64) *Limiter {
	mapMutex.RLock()

	if lim, exists := limiterMap[key]; exists {
		mapMutex.RUnlock()
		return lim
	}

	mapMutex.RUnlock()

	mapMutex.Lock()
	defer mapMutex.Unlock()

	if lim, exists := limiterMap[key]; exists {
		return lim
	}

	lim := NewLimiter(limitType, limitn)
	limiterMap[key] = lim

	return lim
}
