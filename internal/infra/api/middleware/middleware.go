package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     time.Duration
	burst    int
}

type Visitor struct {
	limiter  *time.Ticker
	lastSeen time.Time
	count    int
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) getVisitor(token string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[token]
	if !exists {
		limiter := time.NewTicker(rl.rate)
		rl.visitors[token] = &Visitor{limiter: limiter, lastSeen: time.Now(), count: 0}
		return rl.visitors[token]
	}

	return v
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Second)
		rl.mu.Lock()
		for token, v := range rl.visitors {
			//fmt.Println(v.lastSeen)
			if time.Since(v.lastSeen) > rl.rate {
				delete(rl.visitors, token)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		visitor := rl.getVisitor(token)
		visitor.count++
		if visitor.count > rl.burst {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		select {
		case <-visitor.limiter.C:
			visitor.count = 0
		default:
		}

		next.ServeHTTP(w, r)
	})
}
