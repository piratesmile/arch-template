package middleware

import (
	"arch-template/pkg/errdefs"
	"arch-template/pkg/response"
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	errEmptyIP = errors.New("rate limit: empty ip")
)

type visitor struct {
	bucket    *rate.Limiter
	lastVisit time.Time
}

type Limiter struct {
	mu       sync.RWMutex
	visitors map[string]*visitor
	rate     rate.Limit
	count    int
	ttl      time.Duration
}

func newLimiter(rate rate.Limit, count int, ttl time.Duration) *Limiter {
	return &Limiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		count:    count,
		ttl:      ttl,
	}
}

func (l *Limiter) cleanVisitors() {
	for {
		time.Sleep(time.Minute * 5)

		l.mu.Lock()

		n := time.Now()
		for k, v := range l.visitors {
			if n.Sub(v.lastVisit) > l.ttl {
				delete(l.visitors, k)
			}
		}

		l.mu.Unlock()
	}
}

func (l *Limiter) getVisitor(ip string) *rate.Limiter {

	l.mu.RLock()
	v, ok := l.visitors[ip]
	l.mu.RUnlock()

	if ok {
		v.lastVisit = time.Now()
		return v.bucket
	}

	l.mu.Lock()
	limiter := rate.NewLimiter(l.rate, l.count)
	l.visitors[ip] = &visitor{
		bucket:    limiter,
		lastVisit: time.Now(),
	}
	l.mu.Unlock()
	return limiter
}

func (m *middleware) Limit(rps rate.Limit, count int, ttl time.Duration) gin.HandlerFunc {
	limiter := newLimiter(rps, count, ttl)

	go limiter.cleanVisitors()

	return func(ctx *gin.Context) {
		ip := ctx.RemoteIP()
		if ip == "" {
			response.Error(ctx, errEmptyIP)
			return
		}

		if !limiter.getVisitor(ip).Allow() {
			response.Error(ctx, errdefs.ErrTooManyRequests)
			return
		}
	}

}
